package app

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"time"
)

type App struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *App {
	app := &App{
		logger: logger,
	}
	return app
}

func (a *App) Start(ctx context.Context) error {

	s, err := gocron.NewScheduler()

	if err != nil {
		return err
	}

	_, err = s.NewJob(
		gocron.DurationJob(2*time.Second),
		gocron.NewTask(func() any { return 2 + 2 }))

	if err != nil {
		return err
	}

	s.Start()

	a.logger.Info("server started")

	select {
	case <-ctx.Done():
		return nil
	}

}
