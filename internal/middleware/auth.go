package middleware

import (
	"os"
	"strings"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware() fiber.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-123"
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeUnauthorized,
					Message: "Authorization header required",
				},
			})
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Verify JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInvalidToken,
					Message: "JWT token is invalid or expired",
				},
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInvalidToken,
					Message: "Invalid token claims",
				},
			})
		}

		// Store user info in context for downstream handlers
		c.Locals("userId", claims["userId"])
		c.Locals("playfabId", claims["playfabId"])

		return c.Next()
	}
}
