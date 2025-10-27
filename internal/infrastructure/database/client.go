package database

import (
	"context"
	"fmt"

	"github.com/devgugga/avetis-drive/internal/infrastructure/config"
	"github.com/devgugga/avetis-drive/internal/infrastructure/database/ent"
	"github.com/rs/zerolog"
)

// Client wraps the Ent client with additional functionality
type Client struct {
	*ent.Client
	logger zerolog.Logger
}

// NewClient creates a new database client with Ent
func NewClient(cfg *config.Config, logger zerolog.Logger) (*Client, error) {
	// Create connection
	conn, err := NewConnection(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create Ent client
	entClient := ent.NewClient(ent.Driver(conn.Driver()))

	client := &Client{
		Client: entClient,
		logger: logger,
	}

	return client, nil
}

// AutoMigrate runs automatic migrations
func (c *Client) AutoMigrate(ctx context.Context) error {
	c.logger.Info().Msg("Running database migrations...")

	// Create schema with options to handle if database doesn't exist
	if err := c.Client.Schema.Create(
		ctx,
		// This will create tables if they don't exist, or update them if needed
	); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	c.logger.Info().Msg("Database migrations completed successfully")
	return nil
}

// Ping checks database connectivity
func (c *Client) Ping(ctx context.Context) error {
	// Use a simple query to check connectivity
	_, err := c.Client.User.Query().Limit(1).Count(ctx)
	return err
}

// Close closes the database client
func (c *Client) Close() error {
	c.logger.Info().Msg("Closing database client")
	return c.Client.Close()
}
