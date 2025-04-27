package main

import (
	"context"
	"embed"
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

	dbConn, err := database.Connect(ctx, cfg.DatabaseConfig, migrations, seeds)
	if err != nil {
		log.Panicf("unable to connect to database %v", err)
	}
	defer dbConn.Close()

	db := database.NewTflAlertsDatabase(dbConn)
	twilio := notification.NewTwilioClient(cfg.TwilioConfig)
	smsNotifier, _ := notification.NewSMSNotifier(twilio)
	tflClient, _ := tfl.NewClient(cfg.TflConfig)

	trainSvc := service.LineStatusService{
		TrainsRepo: db.TrainsRepository,
		Tfl:        &tflClient,
	}

	userSvc := service.NotificationService{
		UsersRepo: db.UsersRepository,
		Notifier:  smsNotifier,
	}

	svc := service.DisruptionService{
		TrainService: trainSvc,
		UserService:  userSvc,
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			svc.Start(ctx)
		case <-sigChan:
			log.Println("stopping service...")
			return
		}
	}
}
