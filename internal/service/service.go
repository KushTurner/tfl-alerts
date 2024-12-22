package service

import (
	"context"
	"github.com/kushturner/tfl-alerts/internal/repository"
)

type DisruptionService struct {
	repo repository.Repository
}

func (s DisruptionService) GetUsersToNotify(ctx context.Context, train string) []*repository.User {
	users, _ := s.repo.FindUsersWithDisruptedTrains(ctx, train)

	if len(users) == 0 {
		return nil
	}

	return users
}
