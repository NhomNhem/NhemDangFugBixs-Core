package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware logs request details using structured logging
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Continue stack
		err := c.Next()
		
		// Capture request context
		requestID := c.Locals("requestId")
		userID := c.Locals("userId")
		
		// Calculate latency
		latency := time.Since(start)
		
		status := c.Response().StatusCode()
		
		// Log attributes
		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
			slog.Duration("latency", latency),
		}
		
		if requestID != nil {
			attrs = append(attrs, slog.Any("request_id", requestID))
		}
		
		if userID != nil {
			attrs = append(attrs, slog.Any("user_id", userID))
		}

		msg := "HTTP Request"
		if err != nil {
			msg = err.Error()
			slog.LogAttrs(c.Context(), slog.LevelError, msg, attrs...)
		} else if status >= 400 {
			slog.LogAttrs(c.Context(), slog.LevelWarn, msg, attrs...)
		} else {
			slog.LogAttrs(c.Context(), slog.LevelInfo, msg, attrs...)
		}

		return err
	}
}
