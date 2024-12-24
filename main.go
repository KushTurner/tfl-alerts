package main

import (
	"context"
	"embed"
	"github.com/kushturner/tfl-alerts/internal/api"
	"github.com/kushturner/tfl-alerts/internal/config"
	"github.com/kushturner/tfl-alerts/internal/database"
	"github.com/kushturner/tfl-alerts/internal/notification"
	"github.com/kushturner/tfl-alerts/internal/repository"
	"github.com/kushturner/tfl-alerts/internal/service"
	"log"
	"os"
	"os/signal"
	"time"
)

//go:embed migrations/*.sql
var migrations embed.FS

//go:embed seeds/*.sql
var seeds embed.FS

func main() {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Panicf("unable to load config: %v", err)
	}

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseConfig, migrations, seeds)
	if err != nil {
		log.Panicf("unable to connect to database %v", err)
	}
	defer db.Close()

	repo := repository.NewSQLRepository(db)
	twilio := notification.NewTwilioClient(cfg.TwilioConfig)
	smsNotifier, _ := notification.NewSMSNotifier(twilio)
	tfl, _ := api.NewTflClient(cfg.TflConfig)

	svc := service.NewDisruptionService(repo, smsNotifier, tfl)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			svc.FindUsersAndNotify(ctx)
			svc.PollTrains(ctx)
		case <-sigChan:
			log.Println("stopping service...")
			return
		}
	}
}
