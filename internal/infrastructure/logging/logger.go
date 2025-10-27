package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/devgugga/avetis-drive/internal/infrastructure/config"
	"github.com/rs/zerolog"
)

// NewLogger creates a new zerolog logger with console and file outputs
func NewLogger(cfg *config.Config) (zerolog.Logger, error) {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return zerolog.Logger{}, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Parse log level
	level := parseLogLevel(cfg.Logging.Level)
	zerolog.SetGlobalLevel(level)

	// Console writer - human-readable with colors
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}

	// File output - JSON format
	logFile := filepath.Join(logsDir, "app.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("failed to open log file: %w", err)
	}

	// Multi writer (console + file)
	multi := io.MultiWriter(consoleWriter, file)

	// Create logger
	logger := zerolog.New(multi).
		With().
		Timestamp().
		Logger()

	return logger, nil
}

// parseLogLevel converts string log level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
