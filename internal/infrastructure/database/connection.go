package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/devgugga/avetis-drive/internal/infrastructure/config"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// Connection wraps the database connection
type Connection struct {
	driver *entsql.Driver
	logger zerolog.Logger
}

// NewConnection creates a new database connection
func NewConnection(cfg *config.Config, logger zerolog.Logger) (*Connection, error) {
	// Build connection string
	dsn := cfg.Database.DatabaseURL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
	}

	logger.Debug().Str("dsn", maskPassword(dsn)).Msg("Connecting to database")

	// Try to create database if it doesn't exist
	if err := ensureDatabaseExists(cfg, logger); err != nil {
		logger.Warn().Err(err).Msg("Could not ensure database exists, will try to connect anyway")
	}

	// Open connection with pgx driver
	drv, err := entsql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db := drv.DB()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().
		Str("host", cfg.Database.Host).
		Str("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Msg("Database connection established")

	return &Connection{
		driver: drv,
		logger: logger,
	}, nil
}

// Driver returns the underlying SQL driver for Ent client
func (c *Connection) Driver() *entsql.Driver {
	return c.driver
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.driver != nil {
		c.logger.Info().Msg("Closing database connection")
		return c.driver.Close()
	}
	return nil
}

// ensureDatabaseExists creates the database if it doesn't exist
func ensureDatabaseExists(cfg *config.Config, logger zerolog.Logger) error {
	// Connect to postgres database (which always exists) to create our database
	postgresDB := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := sql.Open("postgres", postgresDB)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, cfg.Database.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if exists {
		logger.Debug().Str("database", cfg.Database.Name).Msg("Database already exists")
		return nil
	}

	// Create database
	logger.Info().Str("database", cfg.Database.Name).Msg("Creating database...")
	createQuery := fmt.Sprintf("CREATE DATABASE %s", quoteDatabaseName(cfg.Database.Name))
	_, err = db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	logger.Info().Str("database", cfg.Database.Name).Msg("Database created successfully")
	return nil
}

// quoteDatabaseName safely quotes a database name for SQL queries
func quoteDatabaseName(name string) string {
	// Replace any existing quotes and wrap in quotes
	name = strings.ReplaceAll(name, `"`, `""`)
	return fmt.Sprintf(`"%s"`, name)
}

// maskPassword masks the password in the DSN for logging
func maskPassword(dsn string) string {
	// Simple masking - in production you might want a more robust solution
	return "postgres://***:***@***"
}
