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
		log.Panicf("unable to load config: %v", err)
	}

	ctx := context.Background()

	store, _ := database.NewDatabase(ctx, cfg.DatabaseConfig)
	defer store.Close()
	usersRepo := store.GetUsersRepository()
	trainsRepo := store.GetTrainsRepository()
	twilio := notification.NewTwilioClient(cfg.TwilioConfig)
	smsNotifier, _ := notification.NewSMSNotifier(twilio)
	tflClient, _ := tfl.NewClient(cfg.TflConfig)

	svc := service.NewDisruptionService(trainsRepo, usersRepo, smsNotifier, tflClient)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			_ = svc.PollTrains(ctx)
			_ = svc.FindUsersAndNotify(ctx)
		case <-sigChan:
			log.Println("stopping service...")
			return
		}
	}
}
