package repository

import (
	"context"
	"github.com/kushturner/tfl-alerts/internal/database"
	"time"
)

type User struct {
	ID           int
	LastNotified time.Time
	Number       string
}

type Repository interface {
	FindUsersWithDisruptedTrains(ctx context.Context, train string) ([]*User, error)
}

type SQLRepository struct {
	db *database.DB
}

func NewSQLRepository(db *database.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}

func (r SQLRepository) FindUsersWithDisruptedTrains(ctx context.Context, train string) ([]*User, error) {

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

	var users []*User

	for rows.Next() {

		user := &User{}
		err := rows.Scan(&user.ID, &user.LastNotified, &user.Number)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}
