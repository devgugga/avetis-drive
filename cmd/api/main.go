package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devgugga/avetis-drive/internal/infrastructure/config"
	apphttp "github.com/devgugga/avetis-drive/internal/infrastructure/http"
	"github.com/devgugga/avetis-drive/internal/infrastructure/logging"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := logging.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info().
		Str("environment", cfg.App.Environment).
		Str("host", cfg.Server.Host).
		Str("port", cfg.Server.Port).
		Msg("Application starting")

	// Initialize HTTP server
	server := apphttp.NewServer(cfg, logger)

	// Start server in a goroutine
	go func() {
		logger.Info().
			Str("address", cfg.Server.Host+":"+cfg.Server.Port).
			Msg("Starting HTTP server")

		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
		os.Exit(1)
	}

	logger.Info().Msg("Server exited gracefully")
}
