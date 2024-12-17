package main

import (
	"github.com/kushturner/tfl-alerts/internal/api"
	"github.com/kushturner/tfl-alerts/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Panicf("unable to load config %v", err)
	}

	client := api.NewTflClient(cfg.TflConfig)

	client.AllCurrentDisruptions()
}
