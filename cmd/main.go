package main

import (
	"fmt"
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

	disruptions, err := client.AllCurrentDisruptions()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(disruptions)
}
