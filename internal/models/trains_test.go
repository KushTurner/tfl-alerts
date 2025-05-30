package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TrainWithSeverity(severity int) *Train {
	return &Train{
		ID:               1,
		Line:             "Jubilee",
		LastUpdated:      time.Now().UTC(),
		PreviousSeverity: 2,
		Severity:         severity,
	}
}

func TestTrain_IsDisrupted(t *testing.T) {
	t.Run("Will return true if a train has severe delays", func(t *testing.T) {
		train := TrainWithSeverity(6)

		assert.True(t, train.IsDisrupted())
	})

	t.Run("Will return true if a train has minor delays", func(t *testing.T) {
		train := TrainWithSeverity(9)

		assert.True(t, train.IsDisrupted())
	})

	t.Run("Will return false if a train does not have severe or minor delays", func(t *testing.T) {
		train := TrainWithSeverity(10)

		assert.False(t, train.IsDisrupted())
	})
}

func TestTrain_SeverityMessage(t *testing.T) {
	t.Run("Will return correct message for 6", func(t *testing.T) {
		train := TrainWithSeverity(6)

		assert.Equal(t, train.SeverityMessage(), "Severe Delays")
	})

	t.Run("Will return correct message for 9", func(t *testing.T) {
		train := TrainWithSeverity(9)

		assert.Equal(t, train.SeverityMessage(), "Minor Delays")
	})
}

func TestTrain_HasSameSeverity(t *testing.T) {
	t.Run("Will return true if previous severity same as current severity", func(t *testing.T) {
		hasSameSeverity := TrainWithSeverity(2).HasSameSeverity()

		assert.True(t, hasSameSeverity)
	})

	t.Run("Will return false if previous severity is not the same as current severity", func(t *testing.T) {
		hasSameSeverity := TrainWithSeverity(9).HasSameSeverity()

		assert.False(t, hasSameSeverity)
	})
}
