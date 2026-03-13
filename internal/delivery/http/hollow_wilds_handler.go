package http

import (
	"fmt"
	"log"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// HollowWildsHandler handles Hollow Wilds game endpoints
type HollowWildsHandler struct {
	authUsecase      usecase.AuthUsecase
	playerUsecase    usecase.PlayerUsecase
	analyticsUsecase usecase.AnalyticsUsecase
}

// NewHollowWildsHandler creates a new Hollow Wilds handler
func NewHollowWildsHandler(
	authUsecase usecase.AuthUsecase,
	playerUsecase usecase.PlayerUsecase,
	analyticsUsecase usecase.AnalyticsUsecase,
) *HollowWildsHandler {
	return &HollowWildsHandler{
		authUsecase:      authUsecase,
		playerUsecase:    playerUsecase,
		analyticsUsecase: analyticsUsecase,
	}
}

// Login handles Hollow Wilds player login
// @Summary Hollow Wilds Login
// @Description Authenticate player with PlayFab ticket and get HW session JWT
// @Tags HollowWilds
// @Accept json
// @Produce json
// @Param X-PlayFab-ID header string true "PlayFab ID"
// @Param request body models.HollowWildsLoginRequest true "Login request"
// @Success 200 {object} models.HollowWildsAuthResponse "Successful login"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /auth/hw/login [post]
func (h *HollowWildsHandler) Login(c *fiber.Ctx) error {
	var req models.HollowWildsLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	overrideID := c.Get("X-PlayFab-ID")
	resp, err := h.authUsecase.Login(c.Context(), req.PlayfabSessionTicket, overrideID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUnauthorized,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(resp)
}

// Refresh handles token refresh
// @Summary Refresh HW Token
// @Description Get a new JWT using a refresh token
// @Tags HollowWilds
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh request"
// @Success 200 {object} models.RefreshTokenResponse "New token"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Invalid refresh token"
// @Router /auth/refresh [post]
func (h *HollowWildsHandler) Refresh(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	resp, err := h.authUsecase.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUnauthorized,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(resp)
}

// Logout handles player logout
// @Summary HW Logout
// @Description Revoke tokens and logout
// @Tags HollowWilds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest false "Optional refresh token to revoke"
// @Success 200 {object} map[string]bool "Success"
// @Router /auth/logout [delete]
func (h *HollowWildsHandler) Logout(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	c.BodyParser(&req) // Optional, don't fail if missing

	// Extract JTI from JWT (in a real app, we'd have a token util)
	// For now, we'll just handle refresh token if provided
	err := h.authUsecase.Logout(c.Context(), req.RefreshToken, "")
	if err != nil {
		log.Printf("Logout warning: %v", err)
	}

	return c.JSON(fiber.Map{"success": true})
}

// GetSave retrieves player save data
// @Summary Load HW Game
// @Description Get the full game state for the player
// @Tags HollowWilds
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.LoadGameResponse "Game save data"
// @Failure 404 {object} models.APIResponse{error=models.APIError} "Save not found"
// @Router /player/save [get]
func (h *HollowWildsHandler) GetSave(c *fiber.Ctx) error {
	playerIDStr := c.Locals("userId").(string)
	playerID, _ := uuid.Parse(playerIDStr)

	save, err := h.playerUsecase.GetSave(c.Context(), playerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve save data",
			},
		})
	}

	if save == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "save_not_found",
				Message: "No save data found for this player",
			},
		})
	}

	return c.JSON(models.LoadGameResponse{
		PlayerID:       save.PlayerID.String(),
		SaveVersion:    save.SaveVersion,
		UpdatedAt:      save.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		World:          save.SaveData.World,
		Player:         save.SaveData.Player,
		Inventory:      save.SaveData.Inventory,
		Sebilah:        save.SaveData.Sebilah,
		Base:           save.SaveData.Base,
		DiscoveredPOIs: save.SaveData.DiscoveredPOIs,
		QuestFlags:     save.SaveData.QuestFlags,
	})
}

// UpdateSave updates player save data
// @Summary Save HW Game
// @Description Persist the game state with version control
// @Tags HollowWilds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param version query int false "Expected current version"
// @Param request body models.SaveGameRequest true "Save data"
// @Success 200 {object} models.SaveGameResponse "Save confirmed"
// @Failure 409 {object} models.VersionConflictError "Version conflict"
// @Router /player/save [put]
func (h *HollowWildsHandler) UpdateSave(c *fiber.Ctx) error {
	playerIDStr := fmt.Sprintf("%v", c.Locals("userId"))
	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Internal server error: auth context invalid",
			},
		})
	}

	var req models.SaveGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	expectedVersion := c.QueryInt("version", 0)
	save, err := h.playerUsecase.SaveGame(c.Context(), playerID, models.GameSaveData{
		World:          req.World,
		Player:         req.Player,
		Inventory:      req.Inventory,
		Sebilah:        req.Sebilah,
		Base:           req.Base,
		DiscoveredPOIs: req.DiscoveredPOIs,
		QuestFlags:     req.QuestFlags,
	}, expectedVersion)

	if err != nil {
		if conflict, ok := err.(*models.VersionConflictError); ok {
			return c.Status(fiber.StatusConflict).JSON(conflict)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(models.SaveGameResponse{
		Success:     true,
		SaveVersion: save.SaveVersion,
		UpdatedAt:   save.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// CreateBackup creates a save backup
// @Summary Create HW Backup
// @Description Manually trigger a save backup
// @Tags HollowWilds
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.BackupResponse "Backup created"
// @Router /player/save/backup [post]
func (h *HollowWildsHandler) CreateBackup(c *fiber.Ctx) error {
	playerIDStr := c.Locals("userId").(string)
	playerID, _ := uuid.Parse(playerIDStr)

	backup, err := h.playerUsecase.CreateBackup(c.Context(), playerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(models.BackupResponse{
		Success:   true,
		BackupID:  backup.ID.String(),
		CreatedAt: backup.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetBackups lists player backups
// @Summary List HW Backups
// @Description Get list of all save backups for the player
// @Tags HollowWilds
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.BackupListResponse "Backup list"
// @Router /player/save/backups [get]
func (h *HollowWildsHandler) GetBackups(c *fiber.Ctx) error {
	playerIDStr := c.Locals("userId").(string)
	playerID, _ := uuid.Parse(playerIDStr)

	backups, err := h.playerUsecase.GetBackups(c.Context(), playerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve backups",
			},
		})
	}

	var backupInfos []models.BackupInfo
	for _, b := range backups {
		backupInfos = append(backupInfos, models.BackupInfo{
			BackupID:    b.ID.String(),
			SaveVersion: b.SaveVersion,
			CreatedAt:   b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return c.JSON(models.BackupListResponse{
		Backups: backupInfos,
	})
}

// RestoreFromBackup restores save data from a backup
// @Summary Restore HW Backup
// @Description Restore game save data from a selected backup
// @Tags HollowWilds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.RestoreRequest true "Restore request"
// @Success 200 {object} models.SaveGameResponse "Restore confirmed"
// @Failure 404 {object} models.APIResponse{error=models.APIError} "Backup not found"
// @Router /player/save/restore [post]
func (h *HollowWildsHandler) RestoreFromBackup(c *fiber.Ctx) error {
	playerIDStr := c.Locals("userId").(string)
	playerID, _ := uuid.Parse(playerIDStr)

	var req models.RestoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	backupID, err := uuid.Parse(req.BackupID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid backup ID format",
			},
		})
	}

	save, err := h.playerUsecase.RestoreFromBackup(c.Context(), playerID, backupID)
	if err != nil {
		if err.Error() == "backup not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "backup_not_found",
					Message: "Backup not found or does not belong to this player",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(models.SaveGameResponse{
		Success:     true,
		SaveVersion: save.SaveVersion,
		UpdatedAt:   save.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// TrackEvents records analytics events
// @Summary Track HW Events
// @Description Submit a batch of analytics events
// @Tags HollowWilds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.AnalyticsEventsRequest true "Analytics events"
// @Success 200 {object} models.AnalyticsEventsResponse "Results"
// @Router /analytics/events [post]
func (h *HollowWildsHandler) TrackEvents(c *fiber.Ctx) error {
	var req models.AnalyticsEventsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	var playerID *uuid.UUID
	if userIDStr, ok := c.Locals("userId").(string); ok {
		id, _ := uuid.Parse(userIDStr)
		playerID = &id
	}

	accepted, rejected, err := h.analyticsUsecase.TrackEvents(c.Context(), playerID, req.Events)
	if err != nil {
		log.Printf("Analytics error: %v", err)
	}

	return c.JSON(models.AnalyticsEventsResponse{
		Accepted: accepted,
		Rejected: rejected,
	})
}
