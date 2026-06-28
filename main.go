package main

import (
	"context"
	"github.com/kushturner/tfl-alerts/internal/config"
	"github.com/kushturner/tfl-alerts/internal/database"
	"github.com/kushturner/tfl-alerts/internal/notification"
	"github.com/kushturner/tfl-alerts/internal/service"
	"github.com/kushturner/tfl-alerts/internal/tfl"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	ctx := context.Background()

	store, err := database.NewDatabase(ctx, cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("unable to initialise database: %v", err)
	}
	defer store.Close()

	usersRepo := store.GetUsersRepository()
	trainsRepo := store.GetTrainsRepository()

	twilio := notification.NewTwilioClient(cfg.TwilioConfig)
	smsNotifier, err := notification.NewSMSNotifier(twilio)
	if err != nil {
		log.Fatalf("unable to initialise SMS notifier: %v", err)
	}

	tflClient, err := tfl.NewClient(cfg.TflConfig)
	if err != nil {
		log.Fatalf("unable to initialise TFL client: %v", err)
	}

	svc := service.NewDisruptionService(trainsRepo, usersRepo, smsNotifier, tflClient)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			if err := svc.PollTrains(ctx); err != nil {
				log.Printf("poll trains: %v", err)
			}
			if err := svc.FindUsersAndNotify(ctx); err != nil {
				log.Printf("find users and notify: %v", err)
			}
		case <-sigChan:
			log.Println("stopping service...")
			return
		}
	}
}
