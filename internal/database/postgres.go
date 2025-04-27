package database

import (
	"context"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

//go:embed seeds/*.sql
var seeds embed.FS

type Config struct {
	ConnStr string
}

func Connect(ctx context.Context, cfg *Config) (*DB, error) {
	pgx, err := pgxpool.New(ctx, cfg.ConnStr)

	if err != nil {
		return nil, fmt.Errorf("config was incorrect: %v", err)
	}

	if err := pgx.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DB{pgx}, nil
}

func (db *DB) RunMigrate() error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, db.Config().ConnString())
	if err != nil {
		return fmt.Errorf("migrate new: %s", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func (db *DB) RunSeed(ctx context.Context) error {
	dirName := "seeds"

	sf, err := seeds.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("failed to read seed files: %w", err)
	}

	for _, f := range sf {
		seed, _ := seeds.ReadFile(dirName + "/" + f.Name())

		if _, err := db.Exec(ctx, string(seed)); err != nil {
			return fmt.Errorf("failed to execute seed: %w", err)
		}
	}

	return nil
}
