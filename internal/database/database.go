package database

import (
	"context"
	"embed"
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

func Connect(ctx context.Context, cfg *Config, migrations fs.FS, seeds embed.FS) (*DB, error) {
	pgx, err := pgxpool.New(ctx, cfg.ConnStr)

	if err != nil {
		return nil, fmt.Errorf("config was incorrect: %v", err)
	}

	if err := pgx.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err := runMigration(migrations, cfg.ConnStr); err != nil {
		return nil, fmt.Errorf("failed to run seed: %w", err)
	}

	if err := runSeed(ctx, pgx, seeds); err != nil {
		return nil, fmt.Errorf("failed to run seed: %w", err)
	}

	return &DB{pgx}, nil
}

func runMigration(m fs.FS, connStr string) error {
	source, err := iofs.New(m, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, connStr)
	if err != nil {
		return fmt.Errorf("migrate new: %s", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func runSeed(ctx context.Context, db *pgxpool.Pool, s embed.FS) error {
	dirName := "seeds"

	sf, err := s.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("failed to read seed files: %w", err)
	}

	for _, f := range sf {
		seed, _ := s.ReadFile(dirName + "/" + f.Name())

		if _, err := db.Exec(ctx, string(seed)); err != nil {
			return fmt.Errorf("failed to execute seed: %w", err)
		}
	}

	return nil
}
