package main

import (
	"embed"
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/config"
	"log"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	cfg, err := config.LoadAppConfig()

	if err != nil {
		log.Panicf("unable to load config: %v", err)
	}

	fmt.Println(cfg)
}
