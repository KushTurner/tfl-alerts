package service

import (
	"context"
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/database"
	"github.com/kushturner/tfl-alerts/internal/notification"
	"github.com/kushturner/tfl-alerts/internal/tfl"
	"log"
)

type DisruptionService struct {
	Database database.TflAlertsDatabase
	Notifier notification.Notifier
	Tfl      tfl.Client
}

func NewDisruptionService(database database.TflAlertsDatabase, notifier notification.Notifier, tfl tfl.Client) DisruptionService {
	return DisruptionService{
		Database: database,
		Notifier: notifier,
		Tfl:      tfl,
	}
}

func (s DisruptionService) FindUsersAndNotify(ctx context.Context) error {
	trains, _ := s.Database.TrainsRepository.FindTrainsThatAreWithinWindow(ctx)

	for _, t := range trains {
		if t.PreviousSeverity == t.Severity {
			continue
		}

		if t.IsDisrupted() {
			users, _ := s.Database.UsersRepository.FindUsersWithDisruptedTrains(ctx, t.Line)
			for _, u := range users {
				msg := fmt.Sprintf("There are %v on the %v.", t.SeverityMessage(), t.Line)

				if err := s.Notifier.Notify(msg, u.Number); err != nil {
					log.Printf("unable to notify user: %v", err)
				}

				if err := s.Database.UsersRepository.UpdateUserLastNotified(ctx, u.ID); err != nil {
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
		err := s.Database.TrainsRepository.UpdateTrainStatus(ctx, trainStatus.Name, trainStatus.LineStatuses[0].StatusSeverity)

		if err != nil {
			log.Printf("unable to update train status: %v", err)
		}
	}

	return nil
}
