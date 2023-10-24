package main

import (
	"database/sql"
	"github.com/superhacker2002/personality-service/internal/config"
	"github.com/superhacker2002/personality-service/internal/infrastructure/postgres/migration"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	configs, err := config.New()
	if err != nil {
		log.Fatalf("config loading failed: %v", err)
	}

	db, err := sql.Open("postgres", configs.Db)
	if err != nil {
		log.Fatalf("failed to open connection with database: %v", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("failed to close connection with database: %v", err)
		}
	}()

	err = postgres.Migrate("file://internal/infrastructure/postgres/migration", db)
	if err != nil {
		log.Fatal(err)
	}

}
