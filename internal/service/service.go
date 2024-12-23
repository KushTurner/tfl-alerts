package service

import (
	"context"
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/notification"
	"github.com/kushturner/tfl-alerts/internal/repository"
)

type DisruptionService struct {
	Repo     repository.Repository
	Notifier notification.Notifier
}

func (s DisruptionService) FindUsersAndNotify(ctx context.Context) error {
	trains, _ := s.Repo.FindTrainsThatAreWithinWindow(ctx)

	if len(trains) == 0 {
		return nil
	}

	for _, t := range trains {
		if t.PreviousSeverity == t.Severity {
			continue
		}

		if t.Severity == 9 || t.Severity == 6 {
			users, _ := s.Repo.FindUsersWithDisruptedTrains(ctx, t.Line)
			for _, u := range users {
				msg := fmt.Sprintf("My glorious King, there are %v on the %v", severity[t.Severity], t.Line)
				s.Notifier.Notify(msg, u.Number)
				s.Repo.UpdateUserLastNotified(ctx, u.ID)
			}
		}

	}
	return nil
}

var severity = map[int]string{
	1:  "Closed",
	2:  "Suspended",
	3:  "Part Suspended",
	5:  "Part Closure",
	6:  "Severe Delays",
	7:  "Reduced service",
	9:  "Minor Delays",
	11: "Part Closed",
	14: "Change of Frequency",
	16: "Not Running",
	17: "Issues Reported",
	20: "Service Closed",
}
