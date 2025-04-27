package database

import (
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

func NewTestDatabase(t *testing.T, connStr string) *DB {
	db, err := NewDatabase(t.Context(), &Config{ConnStr: connStr})
	if err != nil {
		log.Printf("unable to initialize database %s", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return db
}

func CreatePostgresInstance(t *testing.T) string {
	container, err := postgres.Run(t.Context(),
		"postgres:17-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)))
	if err != nil {
		log.Printf("failed to start container: %s", err)
	}

	connStr, _ := container.ConnectionString(t.Context(), "sslmode=disable")
	return connStr
}
