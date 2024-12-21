package main

import (
	"context"
	"embed"
	"github.com/kushturner/tfl-alerts/internal/config"
	"github.com/kushturner/tfl-alerts/internal/database"
	"log"
)

//go:embed migrations/*.sql
var migrations embed.FS

//go:embed seeds/*.sql
var seeds embed.FS

func main() {
	cfg, err := config.LoadAppConfig()
	ctx := context.Background()

	if err != nil {
		log.Panicf("unable to load config: %v", err)
	}

	db, err := database.Connect(ctx, cfg.DatabaseConfig, migrations, seeds)
	if err != nil {
		log.Panicf("unable to connect to database %v", err)
	}

	defer db.Close()
}
