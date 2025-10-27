package main

import (
	"log"

	"github.com/devgugga/avetis-drive/internal/infrastructure/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Application starting in %s environment", cfg.App.Environment)
	log.Printf("Server will listen on %s:%s", cfg.Server.Host, cfg.Server.Port)
}
