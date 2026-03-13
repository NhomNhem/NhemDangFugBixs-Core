package http

import (
	"fmt"
	"log"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TalentHandler handles character talent endpoints
type TalentHandler struct {
	talentUsecase usecase.TalentUsecase
}

// NewTalentHandler creates a new talent handler
func NewTalentHandler(talentUsecase usecase.TalentUsecase) *TalentHandler {
	return &TalentHandler{
		talentUsecase: talentUsecase,
	}
}

// UpgradeTalent handles talent upgrade request
// @Summary Upgrade a talent
// @Description Upgrade a talent using gold currency
// @Tags Talents
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param request body models.TalentUpgradeRequest true "Talent upgrade request"
// @Success 200 {object} models.APIResponse{data=models.TalentUpgradeResponse} "Talent upgraded successfully"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request, insufficient gold, or max level"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /talents/upgrade [post]
// @Security BearerAuth
func (h *TalentHandler) UpgradeTalent(c *fiber.Ctx) error {
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

	if req.TalentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "talentId is required",
			},
		})
	}

	response, err := h.talentUsecase.UpgradeTalent(c.Context(), userID, req.TalentID)
	if err != nil {
		log.Printf("Failed to upgrade talent for user %s: %v", userID, err)
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

		// Use a simpler check for insufficient gold
		if fmt.Sprintf("%.17s", errMsg) == "insufficient gold" {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInsufficientFunds,
					Message: errMsg,
				},
			})
		}

		if fmt.Sprintf("%.10s", errMsg) == "invalid talent ID" {
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
// @Summary Get user talents
// @Description Retrieve all talents and their levels for the authenticated user
// @Tags Talents
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Success 200 {object} models.APIResponse{data=[]models.UserTalent} "User talents retrieved successfully"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /talents [get]
// @Security BearerAuth
func (h *TalentHandler) GetTalents(c *fiber.Ctx) error {
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

	talents, err := h.talentUsecase.GetUserTalents(c.Context(), userID)
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
