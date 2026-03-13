package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique ID to each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get from header if already present (e.g. from Cloudflare/Load Balancer)
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set in context for handlers
		c.Locals("requestId", requestID)

		// Set in response header
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
