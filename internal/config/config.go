package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"log"
)

type Config struct {
	Port string `env:"PORT,default=8080"`
	Db   string `env:"DATABASE_URL,default=postgresql://postgres:2587@localhost:5432/cinema?sslmode=disable"`
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
