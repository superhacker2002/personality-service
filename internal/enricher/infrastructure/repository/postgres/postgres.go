package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/superhacker2002/personality-service/internal/enricher/service"
	"github.com/superhacker2002/personality-service/internal/entity"
	"log"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) Person(id int) (entity.Person, error) {
	var (
		patronymic sql.NullString
		p          entity.Person
	)
	err := r.db.QueryRow(`SELECT name, surname, patronymic, age, gender, nationality
		FROM people
		WHERE id = $1`, id).Scan(&p.Name, &p.Surname, &patronymic,
		&p.Age, &p.Gender, &p.Nationality)
	p.Patronymic = patronymic.String

	if errors.Is(err, sql.ErrNoRows) {
		log.Println(service.ErrPersonNotFound)
		return entity.Person{}, service.ErrPersonNotFound
	}

	if err != nil {
		log.Println(err)
		return entity.Person{}, fmt.Errorf("failed to get person by id: %w", err)
	}

	return p, nil
}

func (r Repository) AllPeople(offset, limit int) ([]entity.Person, error) {
	rows, err := r.db.Query(`SELECT name, surname, patronymic, age, gender, nationality 
		FROM people
		OFFSET $1
		LIMIT $2`, offset, limit)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get all people %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	people, err := r.readPeople(rows)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func (r Repository) PeopleByName(name string, offset int, limit int) ([]entity.Person, error) {
	rows, err := r.db.Query(`SELECT name, surname, patronymic, age, gender, nationality 
		FROM people 
		WHERE name = $1
		OFFSET $2
		LIMIT $3`, name, offset, limit)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get people with name %s: %w", name, err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	people, err := r.readPeople(rows)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func (r Repository) AddPerson(name, surname, patronymic, gender, nationality string, age int) (int, error) {
	var id int
	err := r.db.QueryRow(`INSERT INTO people (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`, name, surname, newNullString(patronymic), age, gender, nationality).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("failed to add person to a database: %w", err)
	}

	return id, nil
}

func (r Repository) DeletePerson(id int) (bool, error) {
	res, err := r.db.Exec("DELETE FROM people WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to delete person: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r Repository) UpdatePerson(id int, name, surname, patronymic string) (bool, error) {
	res, err := r.db.Exec(`UPDATE people 
						SET name = $1, surname = $2, patronymic = $3
						WHERE id = $4`, name, surname, newNullString(patronymic), id)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to update person: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (r Repository) readPeople(rows *sql.Rows) ([]entity.Person, error) {
	var (
		people     []entity.Person
		patronymic sql.NullString
	)

	for rows.Next() {
		var p entity.Person
		if err := rows.Scan(&p.Name, &p.Surname, &patronymic,
			&p.Age, &p.Gender, &p.Nationality); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get person: %w", err)
		}
		p.Patronymic = patronymic.String
		people = append(people, p)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error while iterating over people: %w", err)
	}

	if len(people) == 0 {
		log.Println(service.ErrPeopleNotFound)
		return nil, service.ErrPeopleNotFound
	}

	return people, nil
}
