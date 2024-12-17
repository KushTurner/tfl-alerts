package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func Connect(ctx context.Context, cfg *Config) (*DB, error) {
	parsedConfig, _ := pgxpool.ParseConfig(
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))

	pgx, err := pgxpool.NewWithConfig(ctx, parsedConfig)

	if err != nil {
		return nil, fmt.Errorf("config was incorrect %v", err)
	}

	if err := pgx.Ping(ctx); err != nil {
		return nil, fmt.Errorf("was unable to ping database %v", err)
	}

	return &DB{pgx}, nil
}
