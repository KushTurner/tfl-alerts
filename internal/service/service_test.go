package service

import (
	"context"
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestUsersRepo struct {
}

type TestTrainsRepo struct {
}

type TestNotifier struct {
}

func (t TestNotifier) Notify(msg string, to string) error {
	fmt.Printf("Sent '%s' to %s", msg, to)
	return nil
}

func (t TestTrainsRepo) UpdateTrainStatus(ctx context.Context, train string, severity int) error {
	return nil
}

func (t TestTrainsRepo) FindTrainsThatAreWithinWindow(ctx context.Context) ([]*models.Train, error) {
	var trains []*models.Train
	trains = append(trains, &models.Train{
		ID:               999,
		Line:             "Fake Line",
		LastUpdated:      time.Now().UTC(),
		PreviousSeverity: 2,
		Severity:         9,
	})
	return trains, nil
}

func (t TestUsersRepo) FindUsersWithDisruptedTrains(ctx context.Context, train string) ([]*models.User, error) {
	var users []*models.User
	users = append(users, &models.User{
		ID:           1,
		LastNotified: time.Now(),
		Number:       "number",
	})

	return users, nil
}

func (t TestUsersRepo) UpdateUserLastNotified(ctx context.Context, userID int) error {
	return nil
}

func NewTestDisruptionService() DisruptionService {
	return DisruptionService{
		UsersRepo:  TestUsersRepo{},
		TrainsRepo: TestTrainsRepo{},
		Notifier:   TestNotifier{},
	}
}

func TestName(t *testing.T) {
	t.Run("Can notify users who are disrupted", func(t *testing.T) {
		ds := NewTestDisruptionService()

		err := ds.FindUsersAndNotify(t.Context())

		assert.NoError(t, err)
	})

}
