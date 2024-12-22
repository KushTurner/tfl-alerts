package models

import "time"

type User struct {
	ID           int
	LastNotified time.Time
	Number       string
}
