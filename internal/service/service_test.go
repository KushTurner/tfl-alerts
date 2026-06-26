package service

import (
	"context"
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
	LastMessage string
}

func (t *TestNotifier) Notify(msg string, to string) error {
	t.LastMessage = msg
	return nil
}

func (t TestTrainsRepo) UpdateTrainStatus(ctx context.Context, train string, severity int, summary string) error {
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
		Summary:          "Minor delays between A and B due to a signal failure",
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

func NewTestDisruptionService(notifier *TestNotifier) DisruptionService {
	return DisruptionService{
		UsersRepo:  TestUsersRepo{},
		TrainsRepo: TestTrainsRepo{},
		Notifier:   notifier,
	}
}

func TestFindUsersAndNotify(t *testing.T) {
	t.Run("Can notify users who are disrupted", func(t *testing.T) {
		notifier := &TestNotifier{}
		ds := NewTestDisruptionService(notifier)

		err := ds.FindUsersAndNotify(t.Context())

		assert.NoError(t, err)
	})

	t.Run("Notification message uses line name and disruption summary", func(t *testing.T) {
		notifier := &TestNotifier{}
		ds := NewTestDisruptionService(notifier)

		ds.FindUsersAndNotify(t.Context())

		assert.Equal(t, "Fake Line: Minor delays between A and B due to a signal failure", notifier.LastMessage)
	})

	t.Run("Notification message falls back to severity when summary is empty", func(t *testing.T) {
		notifier := &TestNotifier{}
		ds := NewTestDisruptionService(notifier)
		ds.TrainsRepo = TestTrainsRepoNoSummary{}

		ds.FindUsersAndNotify(t.Context())

		assert.Equal(t, "Fake Line: Minor Delays", notifier.LastMessage)
	})
}

type TestTrainsRepoNoSummary struct{}

func (t TestTrainsRepoNoSummary) UpdateTrainStatus(ctx context.Context, train string, severity int, summary string) error {
	return nil
}

func (t TestTrainsRepoNoSummary) FindTrainsThatAreWithinWindow(ctx context.Context) ([]*models.Train, error) {
	return []*models.Train{
		{
			ID:               999,
			Line:             "Fake Line",
			LastUpdated:      time.Now().UTC(),
			PreviousSeverity: 2,
			Severity:         9,
			Summary:          "",
		},
	}, nil
}
