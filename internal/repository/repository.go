package repository

import (
	"context"
	"github.com/kushturner/tfl-alerts/internal/database"
	"github.com/kushturner/tfl-alerts/internal/models"
	"time"
)

type Repository interface {
	FindUsersWithDisruptedTrains(ctx context.Context, train string) ([]*models.User, error)
	UpdateUserLastNotified(ctx context.Context, userID int) error
	UpdateTrainStatus(ctx context.Context, train string, severity int) error
	FindTrainsThatAreWithinWindow(ctx context.Context) ([]*models.Train, error)
}

type SQLRepository struct {
	db *database.DB
}

func NewSQLRepository(db *database.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}

func (r SQLRepository) FindUsersWithDisruptedTrains(ctx context.Context, train string) ([]*models.User, error) {

	sql := `
		SELECT u.id, u.last_notified, u.phone_number
		FROM users u
        		JOIN notification_windows nw ON u.id = nw.user_id
        		JOIN trains t ON nw.train_id = t.id
		WHERE lower(t.line) = lower($1)
  			AND now() BETWEEN nw.start_time AND nw.end_time`

	rows, err := r.db.Query(ctx, sql, train)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var users []*models.User

	for rows.Next() {

		user := &models.User{}
		err := rows.Scan(&user.ID, &user.LastNotified, &user.Number)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r SQLRepository) UpdateUserLastNotified(ctx context.Context, userID int) error {
	sql := `
        UPDATE users 
        SET last_notified = $1
        WHERE id = $2`

	_, err := r.db.Exec(ctx, sql, time.Now(), userID)
	if err != nil {
		return err
	}

	return nil
}

func (r SQLRepository) UpdateTrainStatus(ctx context.Context, train string, severity int) error {
	sql := `
        UPDATE trains 
        SET previous_severity = severity,
        	severity = $1,
        	last_updated = now()
        WHERE lower(line) = lower($2)`

	_, err := r.db.Exec(ctx, sql, train, severity)
	if err != nil {
		return err
	}

	return nil
}

func (r SQLRepository) FindTrainsThatAreWithinWindow(ctx context.Context) ([]*models.Train, error) {
	sql := `
		SELECT DISTINCT t.id, t.line, t.previous_severity, t.severity
		FROM trains t
			JOIN notification_windows nw ON t.id = nw.train_id
		WHERE now() BETWEEN nw.start_time AND nw.end_time`

	rows, err := r.db.Query(ctx, sql)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var trains []*models.Train

	for rows.Next() {
		train := &models.Train{}
		err := rows.Scan(&train.ID, &train.Line, &train.PreviousSeverity, &train.Severity)

		if err != nil {
			return nil, err
		}

		trains = append(trains, train)
	}

	return trains, nil
}
