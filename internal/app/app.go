package app

import (
	"context"
	"log/slog"
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

	a.logger.Info("server started")

	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}
