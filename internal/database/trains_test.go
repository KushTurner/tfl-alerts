package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgresTrainsRepository_FindTrainsThatAreWithinWindow(t *testing.T) {
	day := int(time.Now().UTC().Weekday())

	t.Run("Can find trains within window", func(t *testing.T) {
		postgres := CreatePostgresInstance(t)
		db := NewTestDatabase(t, postgres)
		defer db.Close()

		currTime := time.Now().UTC()
		startTime := currTime.Add(time.Minute * -5)
		endTime := currTime.Add(time.Minute * 5)

		db.Exec(t.Context(), `INSERT INTO trains (id, line, last_updated, severity, previous_severity, summary) VALUES (999, 'line', now(), 2, 2, 'Minor delays between A and B')`)
		db.Exec(t.Context(), `INSERT INTO users (id, last_notified, phone_number) VALUES (1, now(), 'number')`)
		db.Exec(t.Context(), `INSERT INTO notification_windows (id, user_id, train_id, start_time, end_time, weekday) VALUES (1, 1, 999, $1, $2, $3)`, startTime, endTime, day)

		trains, _ := db.GetTrainsRepository().FindTrainsThatAreWithinWindow(t.Context())

		assert.Equal(t, 1, len(trains))
		assert.Equal(t, "Minor delays between A and B", trains[0].Summary)
	})
}

func TestPostgresTrainsRepository_UpdateTrainStatus(t *testing.T) {
	t.Run("Can update train status with summary", func(t *testing.T) {
		postgres := CreatePostgresInstance(t)
		db := NewTestDatabase(t, postgres)
		defer db.Close()

		db.Exec(t.Context(), `INSERT INTO trains (id, line, last_updated, severity, previous_severity) VALUES (999, 'Elizabeth line', now(), 10, 10)`)

		err := db.GetTrainsRepository().UpdateTrainStatus(t.Context(), "Elizabeth line", 9, "Minor delays between Paddington and Shenfield")

		assert.NoError(t, err)

		var summary string
		db.QueryRow(t.Context(), `SELECT summary FROM trains WHERE id = 999`).Scan(&summary)
		assert.Equal(t, "Minor delays between Paddington and Shenfield", summary)
	})
}
