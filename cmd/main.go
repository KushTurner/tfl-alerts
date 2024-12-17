package main

import (
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadAppConfig()

	if err != nil {
		log.Panicf("unable to load config %v", err)
	}

	fmt.Println(cfg)
}
