package api

import (
	"log"

	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/NhomNhem/GameFeel-Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Login handles user login with PlayFab token
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var req models.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	// Validate required fields
	if req.PlayFabID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "playfabId is required",
			},
		})
	}

	// Get PlayFab session token from header
	sessionToken := c.Get("X-PlayFab-SessionToken")
	
	// Validate PlayFab token (currently skipped - implement later)
	if err := h.authService.ValidatePlayFabToken(sessionToken, req.PlayFabID); err != nil {
		log.Printf("PlayFab token validation skipped: %v", err)
		// Continue anyway for development
	}

	// Get or create user
	user, err := h.authService.GetOrCreateUser(c.Context(), req.PlayFabID, req.DisplayName)
	if err != nil {
		log.Printf("Failed to get/create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to authenticate user",
			},
		})
	}

	// Check if user is banned
	if user.IsBanned {
		return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUserBanned,
				Message: "Your account has been banned",
			},
		})
	}

	// Generate JWT
	jwt, expiresIn, err := h.authService.GenerateJWT(user.ID, user.PlayFabID)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to generate authentication token",
			},
		})
	}

	// Return success response
	return c.JSON(models.APIResponse{
		Success: true,
		Data: models.AuthResponse{
			JWT:       jwt,
			User:      *user,
			ExpiresIn: expiresIn,
		},
	})
}
