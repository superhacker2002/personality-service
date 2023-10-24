package postgres

import (
	"database/sql"
	"log"
)

type PostgresRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p PostgresRepository) AddPerson(name, surname, patronymic, gender, nationality string, age int) (int, error) {
	var id int
	err := p.db.QueryRow(`INSERT INTO people (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`, name, surname, patronymic, age, gender, nationality).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}
