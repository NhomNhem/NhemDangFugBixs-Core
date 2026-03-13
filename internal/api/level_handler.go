package api

import (
	"log"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LevelHandler handles level-related endpoints
type LevelHandler struct {
	levelService *services.LevelService
}

// NewLevelHandler creates a new level handler
func NewLevelHandler() *LevelHandler {
	return &LevelHandler{
		levelService: services.NewLevelService(),
	}
}

// CompleteLevel handles level completion submission
// @Summary Complete a level
// @Description Submit level completion with anti-cheat validation
// @Tags Levels
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param request body models.LevelCompletionRequest true "Level completion data"
// @Success 200 {object} models.APIResponse{data=models.LevelCompletionResponse} "Level completed successfully"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request or cheating detected"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /levels/complete [post]
// @Security BearerAuth
func (h *LevelHandler) CompleteLevel(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUnauthorized,
				Message: "User not authenticated",
			},
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID",
			},
		})
	}

	// Parse request body
	var req models.LevelCompletionRequest
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
	if req.LevelID == "" || req.MapID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "levelId and mapId are required",
			},
		})
	}

	if req.TimeSeconds <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "timeSeconds must be greater than 0",
			},
		})
	}

	// Complete level
	response, err := h.levelService.CompleteLevel(c.Context(), userID, &req)
	if err != nil {
		log.Printf("Failed to complete level for user %s: %v", userID, err)

		// Check if it's a cheating error
		if err.Error() == "invalid completion data" {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeCheatingDetected,
					Message: "Invalid completion data detected",
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to process level completion",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    response,
	})
}
