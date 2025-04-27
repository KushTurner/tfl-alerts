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

func (o DisruptionService) Start(ctx context.Context) {
	err := o.TrainService.Update(ctx)
	if err != nil {
		log.Printf("Error updating training data: %v", err)
	}
	lines, err := o.TrainService.Get(ctx)
	if err != nil {
		log.Printf("Error getting training data: %v", err)
	}
	o.UserService.FindAndNotify(ctx, lines)
}

type LineStatusService struct {
	TrainsRepo database.TrainsRepository
	Tfl        *tfl.Client
}

type NotificationService struct {
	UsersRepo database.UsersRepository
	Notifier  notification.Notifier
}

func (ts LineStatusService) Update(ctx context.Context) error {
	status, _ := ts.Tfl.AllCurrentDisruptions()

	for _, trainStatus := range status {
		err := ts.TrainsRepo.UpdateTrainStatus(ctx, trainStatus.Name, trainStatus.LineStatuses[0].StatusSeverity)

		if err != nil {
			log.Printf("unable to update train status: %v", err)
		}
	}
	return nil
}

func (ts LineStatusService) Get(ctx context.Context) ([]*models.Train, error) {
	return ts.TrainsRepo.FindTrainsThatAreWithinWindow(ctx)
}

func (us NotificationService) FindAndNotify(ctx context.Context, trains []*models.Train) error {
	for _, t := range trains {
		if t.PreviousSeverity == t.Severity {
			continue
		}

		if t.IsDisrupted() {
			users, _ := us.UsersRepo.FindUsersWithDisruptedTrains(ctx, t.Line)
			for _, u := range users {
				msg := fmt.Sprintf("There are %v on the %v.", t.SeverityMessage(), t.Line)

				if err := us.Notifier.Notify(msg, u.Number); err != nil {
					log.Printf("unable to notify user: %v", err)
				}

				if err := us.UsersRepo.UpdateUserLastNotified(ctx, u.ID); err != nil {
					log.Printf("unable to update user: %v", err)
				}
			}
		}
	}
	return nil
}
