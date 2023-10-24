package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"

	"github.com/golang-migrate/migrate/v4"
	dStub "github.com/golang-migrate/migrate/v4/database/postgres"
)

func Migrate(connString string, db *sql.DB) error {
	instance, err := dStub.WithInstance(db, &dStub.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver for the db migration: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(connString, "postgres", instance)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			err = fmt.Errorf("failed to close migration: %w", srcErr)
		} else if dbErr != nil {
			err = fmt.Errorf("failed to close migration: %w", dbErr)
		}
	}()

	if err = m.Up(); errors.Is(err, migrate.ErrNoChange) {
		err = nil
	} else if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	return nil
}
