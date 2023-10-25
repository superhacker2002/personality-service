package postgres

import (
	"errors"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func Migrate(sourceURL string, dbURL string) error {
	m, err := migrate.New(sourceURL, dbURL)
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

	return err
}
