package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devgugga/avetis-drive/internal/infrastructure/config"
	"github.com/devgugga/avetis-drive/internal/infrastructure/database"
	"github.com/devgugga/avetis-drive/internal/infrastructure/http/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// Server wraps the Echo instance and provides server lifecycle management
type Server struct {
	echo     *echo.Echo
	config   *config.Config
	logger   zerolog.Logger
	dbClient *database.Client
}

// NewServer creates and configures a new Echo server instance
func NewServer(cfg *config.Config, logger zerolog.Logger, dbClient *database.Client) *Server {
	e := echo.New()

	// Hide Echo banner
	e.HideBanner = true
	e.HidePort = true

	// Configure server settings
	e.Server.ReadTimeout = 30 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	server := &Server{
		echo:     e,
		config:   cfg,
		logger:   logger,
		dbClient: dbClient,
	}

	// Setup middlewares
	server.setupMiddlewares()

	// Setup routes
	server.setupRoutes()

	return server
}

// setupMiddlewares configures all Echo middlewares
func (s *Server) setupMiddlewares() {
	// Request ID - must be first to be available in all other middlewares
	s.echo.Use(middleware.RequestID())

	// Recover from panics
	s.echo.Use(middleware.Recover())

	// CORS
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // TODO: Configure based on environment
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Zerolog Logger - should be after RequestID to include it in logs
	s.echo.Use(middlewares.ZerologLogger(s.logger))
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// API v1 group
	v1 := s.echo.Group("/api/v1")

	// Health check routes (outside versioned API)
	s.setupHealthRoutes()

	_ = v1 // Avoid unused variable warning for now
}

// setupHealthRoutes configures health check endpoints
func (s *Server) setupHealthRoutes() {
	s.echo.GET("/health", s.healthCheck)
	s.echo.GET("/ready", s.readinessCheck)
}

// healthCheck handles liveness probe
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

// readinessCheck handles readiness probe
func (s *Server) readinessCheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	// Check database connectivity
	if err := s.dbClient.Ping(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Database health check failed")
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"status":   "not ready",
			"database": "unavailable",
			"error":    err.Error(),
			"time":     time.Now().UTC(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "ready",
		"database": "connected",
		"time":     time.Now().UTC(),
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	return s.echo.Start(address)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

// Echo returns the underlying Echo instance for advanced configuration if needed
func (s *Server) Echo() *echo.Echo {
	return s.echo
}
