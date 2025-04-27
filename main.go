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

	db := database.NewTflAlertsDatabase(ctx, cfg.DatabaseConfig)
	twilio := notification.NewTwilioClient(cfg.TwilioConfig)
	smsNotifier, _ := notification.NewSMSNotifier(twilio)
	tflClient, _ := tfl.NewClient(cfg.TflConfig)

	svc := service.NewDisruptionService(db, smsNotifier, tflClient)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			svc.PollTrains(ctx)
			svc.FindUsersAndNotify(ctx)
		case <-sigChan:
			log.Println("stopping service...")
			return
		}
	}
}
