package config

import (
	"github.com/kushturner/tfl-alerts/internal/api"
	"github.com/kushturner/tfl-alerts/internal/database"
	"log"
	"os"
)

type AppConfig struct {
	TflConfig      *api.TflConfig
	DatabaseConfig *database.Config
}

func initDatabase() (*database.Config, error) {
	return &database.Config{
		ConnStr: getEnv("DB_URL", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
	}, nil
}

func initTfl() (*api.TflConfig, error) {
	return &api.TflConfig{
		Url: getEnv("TFL_URL", "https://api.tfl.gov.uk"),
	}, nil
}

func LoadAppConfig() (*AppConfig, error) {
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("unable to initialize database config: %v", err)
	}

	tfl, err := initTfl()
	if err != nil {
		log.Fatalf("unable to initialize tfl client config: %v", err)
	}

	return &AppConfig{
		tfl,
		db,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
