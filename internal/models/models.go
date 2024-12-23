package models

import (
	"time"
)

type User struct {
	ID           int
	LastNotified time.Time
	Number       string
}

type Train struct {
	ID               int
	Line             string
	LastUpdated      time.Time
	PreviousSeverity int
	Severity         int
}
