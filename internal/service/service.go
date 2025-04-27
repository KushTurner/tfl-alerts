package service

import (
	"context"
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/database"
	"github.com/kushturner/tfl-alerts/internal/models"
	"github.com/kushturner/tfl-alerts/internal/notification"
	"github.com/kushturner/tfl-alerts/internal/tfl"
	"log"
)

type DisruptionService struct {
	TrainService LineStatusService
	UserService  NotificationService
}

func (ds DisruptionService) Start(ctx context.Context) {
	err := ds.TrainService.Update(ctx)
	if err != nil {
		log.Printf("Error updating training data: %v", err)
	}
	lines, err := ds.TrainService.Get(ctx)
	if err != nil {
		log.Printf("Error getting training data: %v", err)
	}
	ds.UserService.FindAndNotify(ctx, lines)
}

type LineStatusService struct {
	TrainsRepo database.TrainsRepository
	Tfl        *tfl.Client
}

type NotificationService struct {
	UsersRepo database.UsersRepository
	Notifier  notification.Notifier
}

func (ls LineStatusService) Update(ctx context.Context) error {
	status, _ := ls.Tfl.AllCurrentDisruptions()

	for _, trainStatus := range status {
		err := ls.TrainsRepo.UpdateTrainStatus(ctx, trainStatus.Name, trainStatus.LineStatuses[0].StatusSeverity)

		if err != nil {
			log.Printf("unable to update train status: %v", err)
		}
	}
	return nil
}

func (ls LineStatusService) Get(ctx context.Context) ([]*models.Train, error) {
	return ls.TrainsRepo.FindTrainsThatAreWithinWindow(ctx)
}

func (ns NotificationService) FindAndNotify(ctx context.Context, trains []*models.Train) error {
	for _, t := range trains {
		if t.PreviousSeverity == t.Severity {
			continue
		}

		if t.IsDisrupted() {
			users, _ := ns.UsersRepo.FindUsersWithDisruptedTrains(ctx, t.Line)
			for _, u := range users {
				msg := fmt.Sprintf("There are %v on the %v.", t.SeverityMessage(), t.Line)

				if err := ns.Notifier.Notify(msg, u.Number); err != nil {
					log.Printf("unable to notify user: %v", err)
				}

				if err := ns.UsersRepo.UpdateUserLastNotified(ctx, u.ID); err != nil {
					log.Printf("unable to update user: %v", err)
				}
			}
		}
	}
	return nil
}
