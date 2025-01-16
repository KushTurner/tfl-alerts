package service

import (
	"context"
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/api"
	"github.com/kushturner/tfl-alerts/internal/notification"
	"github.com/kushturner/tfl-alerts/internal/repository"
	"log"
)

type DisruptionService struct {
	Repo     repository.Repository
	Notifier notification.Notifier
	Tfl      api.TflClient
}

func NewDisruptionService(repo repository.Repository, notifier notification.Notifier, tfl api.TflClient) DisruptionService {
	return DisruptionService{
		Repo:     repo,
		Notifier: notifier,
		Tfl:      tfl,
	}
}

func (s DisruptionService) FindUsersAndNotify(ctx context.Context) error {
	trains, err := s.Repo.FindTrainsThatAreWithinWindow(ctx)

	if err != nil {
		log.Printf("unable to find users that are within window: %v", err)
	}

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
				msg := fmt.Sprintf("There are %v on the %v.", severity[t.Severity], t.Line)

				if err := s.Notifier.Notify(msg, u.Number); err != nil {
					log.Printf("unable to notify user: %v", err)
				}

				if err := s.Repo.UpdateUserLastNotified(ctx, u.ID); err != nil {
					log.Printf("unable to update user: %v", err)
				}
			}
		}
	}
	return nil
}

func (s DisruptionService) PollTrains(ctx context.Context) error {
	status, err := s.Tfl.AllCurrentDisruptions()

	if err != nil {
		log.Printf("unable to get all current disruptions: %v", err)
	}

	for _, trainStatus := range status {
		err := s.Repo.UpdateTrainStatus(ctx, trainStatus.Name, trainStatus.LineStatuses[0].StatusSeverity)

		if err != nil {
			log.Printf("unable to update train status: %v", err)
		}
	}

	return nil
}

var severity = map[int]string{
	1:  "Closed",
	2:  "Suspended",
	3:  "Part Suspended",
	4:  "Planned Closure",
	5:  "Part Closure",
	6:  "Severe Delays",
	7:  "Reduced service",
	8:  "Bus Service",
	9:  "Minor Delays",
	10: "Good Service",
	11: "Part Closed",
	12: "Exit Only",
	13: "No Step Free Access",
	14: "Change of Frequency",
	15: "Diverted",
	16: "Not Running",
	17: "Issues Reported",
	18: "No Issues",
	19: "Information",
	20: "Service Closed",
}
