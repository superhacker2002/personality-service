package postgres

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (p Repository) AddPerson(name, surname, patronymic, gender, nationality string, age int) (int, error) {
	var id int
	err := p.db.QueryRow(`INSERT INTO people (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`, name, surname, NewNullString(patronymic), age, gender, nationality).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add person to a database: %w", err)
	}

	return id, nil
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
