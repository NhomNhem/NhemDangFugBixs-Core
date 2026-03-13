package http

import (
	"log"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Login handles user login with PlayFab token (Legacy)
// @Summary Login with PlayFab
// @Description Authenticate user with PlayFab session token and get JWT for API access
// @Tags Authentication
// @Accept json
// @Produce json
// @Param X-PlayFab-SessionToken header string true "PlayFab session token"
// @Param request body models.AuthRequest true "Login request with PlayFab ID"
// @Success 200 {object} models.APIResponse{data=models.AuthResponse} "Successful login with JWT token"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Invalid PlayFab token"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /auth/login [post]
// @Security PlayFabToken
func (h *AuthHandler) Login(c *fiber.Ctx) error {
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

	if req.PlayFabID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "playfabId is required",
			},
		})
	}

	sessionToken := c.Get("X-PlayFab-SessionToken")
	displayName := ""
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}

	resp, err := h.authUsecase.LegacyLogin(c.Context(), req.PlayFabID, displayName, sessionToken)
	if err != nil {
		log.Printf("Legacy login failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUnauthorized,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    resp,
	})
}
