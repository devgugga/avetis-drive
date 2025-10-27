package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Load loads configuration from environment variables
// It will attempt to load a .env file if present, but will not fail if it doesn't exist
func Load() (*Config, error) {
	// Load .env file if it exists (optional in production)
	_ = godotenv.Load()

	config := &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5432"),
			Name:        getEnv("DB_NAME", "avetis_drive"),
			User:        getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", ""),
			DatabaseURL: getEnv("DATABASE_URL", ""),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validate required configuration
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// validate ensures required configuration values are present
func validate(cfg *Config) error {
	if cfg.Database.Password == "" && cfg.Database.DatabaseURL == "" {
		return fmt.Errorf("either DB_PASSWORD or DATABASE_URL must be set")
	}
	return nil
}
