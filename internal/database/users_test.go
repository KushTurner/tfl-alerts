package database

import (
	"github.com/stretchr/testify/assert"
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

func createTestContainer(t *testing.T) string {
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

func TestPostgresUsersRepository_FindUsersWithDisruptedTrains(t *testing.T) {
	day := int(time.Now().UTC().Weekday())

	t.Run("Can find users that have trains that are considered disrupted", func(t *testing.T) {
		connStr := createTestContainer(t)
		db := NewTestDatabase(t, connStr)

		currTime := time.Now().UTC()
		startTime := currTime.Add(time.Minute * -10)
		endTime := currTime.Add(time.Minute * 10)
		db.Exec(t.Context(), `INSERT INTO users (id, last_notified, phone_number) VALUES (1, now(), $1)`, "number")
		db.Exec(t.Context(), `INSERT INTO notification_windows (id, user_id, train_id, start_time, end_time, weekday) VALUES (1, 1, 1, $1, $2, $3)`, startTime, endTime, day)

		actual, _ := db.GetUsersRepository().FindUsersWithDisruptedTrains(t.Context(), "Avanti West Coast")

		assert.Equal(t, 1, actual[0].ID)
	})

	t.Run("Can not find users when no users within window", func(t *testing.T) {
		connStr := createTestContainer(t)
		db := NewTestDatabase(t, connStr)

		currTime := time.Now().UTC()
		startTime := currTime.Add(time.Minute * -10)
		endTime := currTime.Add(time.Minute * -5)
		db.Exec(t.Context(), `INSERT INTO users (id, last_notified, phone_number) VALUES (1, now(), $1)`, "number")
		db.Exec(t.Context(), `INSERT INTO notification_windows (id, user_id, train_id, start_time, end_time, weekday) VALUES (1, 1, 1, $1, $2, $3)`, startTime, endTime, day)

		actual, _ := db.GetUsersRepository().FindUsersWithDisruptedTrains(t.Context(), "Avanti West Coast")

		assert.Equal(t, 0, len(actual))
	})
}
