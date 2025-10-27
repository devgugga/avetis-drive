package config

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Environment string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host        string
	Port        string
	Name        string
	User        string
	Password    string
	DatabaseURL string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}
