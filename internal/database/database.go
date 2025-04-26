package database

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

type TflAlertsDatabase struct {
	TrainsRepository
	UsersRepository
}

func NewTflAlertsDatabase(db *DB) TflAlertsDatabase {
	return TflAlertsDatabase{
		TrainsRepository: PostgresTrainsRepository{db},
		UsersRepository:  PostgresUsersRepository{db},
	}
}
