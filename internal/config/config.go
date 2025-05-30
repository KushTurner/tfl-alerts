package config

import (
	"github.com/kushturner/tfl-alerts/internal/database"
	"github.com/kushturner/tfl-alerts/internal/notification"
	"github.com/kushturner/tfl-alerts/internal/tfl"
	"log"
	"os"
)

type AppConfig struct {
	TflConfig      *tfl.Config
	DatabaseConfig *database.Config
	TwilioConfig   *notification.TwilioConfig
}

func initDatabase() (*database.Config, error) {
	return &database.Config{
		ConnStr: getEnv("DB_URL", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
	}, nil
}

func initTfl() (*tfl.Config, error) {
	return &tfl.Config{
		Url: getEnv("TFL_URL", "https://api.tfl.gov.uk"),
	}, nil
}

func initTwilio() (*notification.TwilioConfig, error) {
	return &notification.TwilioConfig{
		From:       getEnv("TWILIO_FROM", ""),
		AccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		AuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
	}, nil
}

func LoadAppConfig() (*AppConfig, error) {
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("unable to initialize database config: %v", err)
	}

	tflCfg, err := initTfl()
	if err != nil {
		log.Fatalf("unable to initialize tflCfg client config: %v", err)
	}

	twilio, err := initTwilio()
	if err != nil {
		log.Fatalf("unable to initialize twilio client config: %v", err)
	}

	return &AppConfig{
		tflCfg,
		db,
		twilio,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
