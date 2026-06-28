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
	TrainsRepo database.TrainsRepository
	UsersRepo  database.UsersRepository
	Notifier   notification.Notifier
	Tfl        tfl.Client
}

func NewDisruptionService(trainsRepo database.TrainsRepository, usersRepo database.UsersRepository, notifier notification.Notifier, tfl tfl.Client) DisruptionService {
	return DisruptionService{
		TrainsRepo: trainsRepo,
		UsersRepo:  usersRepo,
		Notifier:   notifier,
		Tfl:        tfl,
	}
}

func (s DisruptionService) FindUsersAndNotify(ctx context.Context) error {
	trains, err := s.TrainsRepo.FindTrainsThatAreWithinWindow(ctx)
	if err != nil {
		return fmt.Errorf("finding trains within window: %w", err)
	}

	for _, t := range trains {
		if t.HasSameSeverity() {
			continue
		}

		if t.IsDisrupted() {
			users, err := s.UsersRepo.FindUsersWithDisruptedTrains(ctx, t.Line)
			if err != nil {
				return fmt.Errorf("finding users for train %s: %w", t.Line, err)
			}

			for _, u := range users {
				msg := t.NotificationMessage()

				if err := s.Notifier.Notify(msg, u.Number); err != nil {
					log.Printf("unable to notify user %d: %v", u.ID, err)
				}

				if err := s.UsersRepo.UpdateUserLastNotified(ctx, u.ID); err != nil {
					log.Printf("unable to update last notified for user %d: %v", u.ID, err)
				}
			}
		}
	}
	return nil
}

func (s DisruptionService) PollTrains(ctx context.Context) error {
	status, err := s.Tfl.AllCurrentDisruptions()
	if err != nil {
		return fmt.Errorf("fetching disruptions: %w", err)
	}

	for _, trainStatus := range status {
		summary := trainStatus.LineStatuses[0].Disruption.Description
		if err := s.TrainsRepo.UpdateTrainStatus(ctx, trainStatus.Name, trainStatus.LineStatuses[0].StatusSeverity, summary); err != nil {
			log.Printf("unable to update train status for %s: %v", trainStatus.Name, err)
		}
	}

	return nil
}
