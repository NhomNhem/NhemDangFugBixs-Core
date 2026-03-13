package middleware

import (
	"context"
	"log"

	"github.com/NhomNhem/HollowWilds-Backend/internal/database"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AdminMiddleware checks if user is authenticated AND has admin role
// This middleware should be applied AFTER AuthMiddleware
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First check if JWT auth was successful (userID should be in locals)
		userIDStr, ok := c.Locals("userId").(string)
		if !ok || userIDStr == "" {
			log.Printf("Admin middleware: userId not found in context")
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeUnauthorized,
					Message: "Authentication required",
				},
			})
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			log.Printf("Admin middleware: invalid userId format: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInvalidRequest,
					Message: "Invalid user ID format",
				},
			})
		}

		// Check if user is admin
		db := database.GetDB()
		if db == nil {
			log.Printf("Admin middleware: database not available")
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInternalError,
					Message: "Database connection error",
				},
			})
		}

		var isAdmin bool
		err = db.QueryRow(context.Background(),
			"SELECT is_admin FROM users WHERE id = $1",
			userID,
		).Scan(&isAdmin)

		if err != nil {
			log.Printf("Admin middleware: failed to check admin status for user %s: %v", userID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInternalError,
					Message: "Failed to verify admin status",
				},
			})
		}

		if !isAdmin {
			log.Printf("Admin middleware: access denied for non-admin user %s", userID)
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "FORBIDDEN",
					Message: "Admin access required. You do not have permission to access this resource.",
				},
			})
		}

		// User is admin, store in context and continue
		c.Locals("isAdmin", true)
		log.Printf("Admin middleware: admin access granted for user %s", userID)
		return c.Next()
	}
}
