package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgresUsersRepository_FindUsersWithDisruptedTrains(t *testing.T) {
	day := int(time.Now().UTC().Weekday())

	t.Run("Can find users that have trains that are considered disrupted", func(t *testing.T) {
		connStr := CreatePostgresInstance(t)
		db := NewTestDatabase(t, connStr)

		currTime := time.Now().UTC()
		startTime := currTime.Add(time.Minute * -10)
		endTime := currTime.Add(time.Minute * 10)
		db.Exec(t.Context(), `INSERT INTO users (id, last_notified, phone_number) VALUES (1, now(), 'number')`)
		db.Exec(t.Context(), `INSERT INTO notification_windows (id, user_id, train_id, start_time, end_time, weekday) VALUES (1, 1, 1, $1, $2, $3)`, startTime, endTime, day)

		actual, _ := db.GetUsersRepository().FindUsersWithDisruptedTrains(t.Context(), "Avanti West Coast")

		assert.Equal(t, 1, actual[0].ID)
	})

	t.Run("Can not find users when no users within window", func(t *testing.T) {
		connStr := CreatePostgresInstance(t)
		db := NewTestDatabase(t, connStr)

		currTime := time.Now().UTC()
		startTime := currTime.Add(time.Minute * -5)
		endTime := currTime.Add(time.Minute * -5)
		db.Exec(t.Context(), `INSERT INTO users (id, last_notified, phone_number) VALUES (1, now(), 'number')`)
		db.Exec(t.Context(), `INSERT INTO notification_windows (id, user_id, train_id, start_time, end_time, weekday) VALUES (1, 1, 1, $1, $2, $3)`, startTime, endTime, day)

		actual, _ := db.GetUsersRepository().FindUsersWithDisruptedTrains(t.Context(), "Avanti West Coast")

		assert.Equal(t, 0, len(actual))
	})
}
