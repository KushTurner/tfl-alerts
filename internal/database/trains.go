package database

import (
	"context"
	"github.com/kushturner/tfl-alerts/internal/models"
	"time"
)

type TrainsRepository interface {
	UpdateTrainStatus(ctx context.Context, train string, severity int, summary string) error
	FindTrainsThatAreWithinWindow(ctx context.Context) ([]*models.Train, error)
}

type PostgresTrainsRepository struct {
	Db *DB
}

func (r PostgresTrainsRepository) UpdateTrainStatus(ctx context.Context, train string, severity int, summary string) error {
	sql := `
        UPDATE trains
        SET previous_severity = severity,
        	severity = $1,
        	summary = $2,
        	last_updated = now()
        WHERE lower(line) = lower($3)`

	_, err := r.Db.Exec(ctx, sql, severity, summary, train)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresTrainsRepository) FindTrainsThatAreWithinWindow(ctx context.Context) ([]*models.Train, error) {
	sql := `
		SELECT DISTINCT t.id, t.line, t.previous_severity, t.severity, COALESCE(t.summary, '')
		FROM trains t
			JOIN notification_windows nw ON t.id = nw.train_id
		WHERE nw.weekday = $1
		    AND CURRENT_TIME BETWEEN nw.start_time AND nw.end_time`

	weekday := int(time.Now().UTC().Weekday())

	rows, err := r.Db.Query(ctx, sql, weekday)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var trains []*models.Train

	for rows.Next() {
		train := &models.Train{}
		err := rows.Scan(&train.ID, &train.Line, &train.PreviousSeverity, &train.Severity, &train.Summary)

		if err != nil {
			return nil, err
		}

		trains = append(trains, train)
	}

	return trains, nil
}
