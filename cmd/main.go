package main

import (
	"database/sql"
	"github.com/superhacker2002/personality-service/internal/config"
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

}
