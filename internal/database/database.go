package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"io/fs"
)

type DB struct {
	*pgxpool.Pool
}

type Config struct {
	ConnStr string
}

func Connect(ctx context.Context, cfg *Config, migrations fs.FS) (*DB, error) {
	pgx, err := pgxpool.New(ctx, cfg.ConnStr)

	if err != nil {
		return nil, fmt.Errorf("config was incorrect: %v", err)
	}

	if err := pgx.Ping(ctx); err != nil {
		return nil, fmt.Errorf("was unable to ping database: %v", err)
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("migrate new: %s", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DB{pgx}, nil
}
