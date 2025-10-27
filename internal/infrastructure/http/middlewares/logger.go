package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// ZerologLogger creates a middleware that logs HTTP requests using zerolog
func ZerologLogger(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Get request ID from context (set by RequestID middleware)
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()
			latency := time.Since(start)

			// Create log event
			event := logger.Info()
			if err != nil {
				event = logger.Error().Err(err)
			}

			event.
				Str("request_id", requestID).
				Str("remote_ip", c.RealIP()).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Dur("latency", latency).
				Msg("HTTP request")

			return nil
		}
	}
}
