package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/superhacker2002/personality-service/internal/config"
	"github.com/superhacker2002/personality-service/internal/enricher/controller/httphandler"
	"github.com/superhacker2002/personality-service/internal/enricher/infrastructure/postgres/migration"
	"github.com/superhacker2002/personality-service/internal/enricher/infrastructure/repository/postgres"
	"github.com/superhacker2002/personality-service/internal/enricher/service"
	"log"
	"net/http"
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

	err = postgres.Migrate("file://internal/enricher/infrastructure/postgres/migration/migrations", configs.Db)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	repo := repository.New(db)
	serv := service.New(repo)
	httphandler.New(serv).SetRoutes(router)

	log.Fatal(http.ListenAndServe(":"+configs.Port, router))
}
