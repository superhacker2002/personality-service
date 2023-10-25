package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
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

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
