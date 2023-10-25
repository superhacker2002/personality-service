package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"log"
)

type Config struct {
	Port          string `env:"PORT,default=8080"`
	Db            string `env:"DATABASE_URL,default=postgresql://postgres:postgres@localhost:5432/people?sslmode=disable"`
	MigrationPath string `env:"MIGRATION_FILES_PATH, default=file://internal/enricher/infrastructure/postgres/migration/migrations"`
}

func New() (Config, error) {
	var c Config
	if err := godotenv.Load(); err != nil {
		log.Println("config loading from .env failed:", err)
	}

	ctx := context.Background()
	if err := envconfig.Process(ctx, &c); err != nil {
		return c, err
	}

	return c, nil
}
