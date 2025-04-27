package database

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DB struct {
	*pgxpool.Pool
}

type TflAlertsDatabase struct {
	TrainsRepository
	UsersRepository
}

func NewTflAlertsDatabase(ctx context.Context, config *Config) TflAlertsDatabase {
	db, _ := Connect(ctx, config)

	err := RunMigrate(config.ConnStr)
	if err != nil {
		log.Panicf("unable to run migration %v", err)
	}

	err = RunSeed(ctx, db.Pool)
	if err != nil {
		log.Panicf("unable to run seed %v", err)
	}

	return TflAlertsDatabase{
		TrainsRepository: PostgresTrainsRepository{db},
		UsersRepository:  PostgresUsersRepository{db},
	}
}
