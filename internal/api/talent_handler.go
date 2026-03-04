package api

import (
	"log"

	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/NhomNhem/GameFeel-Backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TalentHandler handles talent-related endpoints
type TalentHandler struct {
	talentService *services.TalentService
}

// NewTalentHandler creates a new talent handler
func NewTalentHandler() *TalentHandler {
	return &TalentHandler{
		talentService: services.NewTalentService(),
	}
}

// UpgradeTalent handles talent upgrade request
// POST /api/v1/talents/upgrade
func (h *TalentHandler) UpgradeTalent(c *fiber.Ctx) error {
	// Get user ID from context
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
	var req models.TalentUpgradeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	// Validate talent ID
	if req.TalentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "talentId is required",
			},
		})
	}

	// Upgrade talent
	response, err := h.talentService.UpgradeTalent(c.Context(), userID, req.TalentID)
	if err != nil {
		log.Printf("Failed to upgrade talent for user %s: %v", userID, err)
		
		// Check error type
		errMsg := err.Error()
		if errMsg == "talent already at max level" {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInvalidRequest,
					Message: "Talent is already at maximum level",
				},
			})
		}
		
		if len(errMsg) > 17 && errMsg[:17] == "insufficient gold" {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInsufficientFunds,
					Message: errMsg,
				},
			})
		}
		
		if len(errMsg) > 10 && errMsg[:10] == "invalid talent ID" {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInvalidRequest,
					Message: "Invalid talent ID",
				},
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to upgrade talent",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetTalents gets all talents for the authenticated user
// GET /api/v1/talents
func (h *TalentHandler) GetTalents(c *fiber.Ctx) error {
	// Get user ID from context
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

	// Get talents
	talents, err := h.talentService.GetUserTalents(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get talents for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve talents",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    talents,
	})
}
