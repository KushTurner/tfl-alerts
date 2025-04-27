package database

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func NewDatabase(ctx context.Context, config *Config) (*DB, error) {
	db, err := Connect(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize postgres: %v", err)
	}
	err = db.RunMigrate()
	if err != nil {
		return nil, fmt.Errorf("unable to run migration: %v", err)
	}
	err = db.RunSeed(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to run seed: %v", err)
	}

	return db, nil
}

func (db *DB) GetUsersRepository() UsersRepository {
	return PostgresUsersRepository{db}
}

func (db *DB) GetTrainsRepository() TrainsRepository {
	return PostgresTrainsRepository{db}
}
