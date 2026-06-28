package database

import (
	"context"
	"github.com/kushturner/tfl-alerts/internal/models"
	"time"
)

type UsersRepository interface {
	FindUsersWithDisruptedTrains(ctx context.Context, train string) ([]*models.User, error)
	UpdateUserLastNotified(ctx context.Context, userID int) error
}

type PostgresUsersRepository struct {
	Db *DB
}

var londonLocation, _ = time.LoadLocation("Europe/London")

func (r PostgresUsersRepository) FindUsersWithDisruptedTrains(ctx context.Context, train string) ([]*models.User, error) {

	sql := `
		SELECT u.id, u.last_notified, u.phone_number
		FROM users u
        		JOIN notification_windows nw ON u.id = nw.user_id
        		JOIN trains t ON nw.train_id = t.id
		WHERE lower(t.line) = lower($1)
		  	AND $2 = nw.weekday
  			AND $3::time BETWEEN nw.start_time AND nw.end_time`

	now := time.Now().In(londonLocation)
	weekday := int(now.Weekday())
	currentTime := now.Format("15:04:05")

	rows, err := r.Db.Query(ctx, sql, train, weekday, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r PostgresUsersRepository) UpdateUserLastNotified(ctx context.Context, userID int) error {
	sql := `
        UPDATE users 
        SET last_notified = $1
        WHERE id = $2`

	_, err := r.Db.Exec(ctx, sql, time.Now().UTC(), userID)
	if err != nil {
		return err
	}

	return nil
}
