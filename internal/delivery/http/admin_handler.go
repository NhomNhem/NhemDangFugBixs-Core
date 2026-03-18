package http

import (
	"log"
	"strconv"
	"strings"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AdminHandler handles administrative endpoints
type AdminHandler struct {
	adminUsecase usecase.AdminUsecase
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminUsecase usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}

// SearchUsers handles user search request
// @Summary Search users
// @Description Search users by PlayFab ID, email, or username (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param q query string true "Search query" minlength(3)
// @Param page query int false "Page number" default(1)
// @Param perPage query int false "Results per page (max 100)" default(20)
// @Success 200 {object} models.APIResponse{data=models.UserSearchResponse} "Search results"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/users/search [get]
// @Security BearerAuth
func (h *AdminHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if len(query) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Search query must be at least 3 characters",
			},
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.Query("perPage", "20"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	result, err := h.adminUsecase.SearchUsers(c.Context(), query, page, perPage)
	if err != nil {
		log.Printf("Failed to search users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to search users",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetUserProfile handles user profile request
// @Summary Get user profile
// @Description Get detailed profile for a specific user (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} models.APIResponse{data=models.UserProfile} "User profile"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid user ID"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 404 {object} models.APIResponse{error=models.APIError} "User not found"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/users/{userId}/profile [get]
// @Security BearerAuth
func (h *AdminHandler) GetUserProfile(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID format",
			},
		})
	}

	profile, err := h.adminUsecase.GetProfile(c.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "NOT_FOUND",
					Message: "User not found",
				},
			})
		}

		log.Printf("Failed to get user profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to get user profile",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    profile,
	})
}

// AdjustGold handles manual gold adjustment
// @Summary Adjust user gold
// @Description Manually adjust user's gold balance (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param userId path string true "User ID (UUID)"
// @Param request body models.AdjustGoldRequest true "Gold adjustment request"
// @Success 200 {object} models.APIResponse{data=models.AdjustGoldResponse} "Gold adjusted"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/users/{userId}/adjust-gold [post]
// @Security BearerAuth
func (h *AdminHandler) AdjustGold(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)

	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID format",
			},
		})
	}

	var req models.AdjustGoldRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	if req.Amount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Amount cannot be zero",
			},
		})
	}

	if len(req.Reason) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Reason must be at least 10 characters",
			},
		})
	}

	ipAddress := c.IP()
	result, err := h.adminUsecase.AdjustGold(c.Context(), adminID, userID, req.Amount, req.Reason, ipAddress)
	if err != nil {
		log.Printf("Failed to adjust gold: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// BanUser handles user ban request
// @Summary Ban user
// @Description Ban a user with reason and optional duration (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param userId path string true "User ID (UUID)"
// @Param request body models.BanUserRequest true "Ban request"
// @Success 200 {object} models.APIResponse{data=models.BanUserResponse} "User banned"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/users/{userId}/ban [post]
// @Security BearerAuth
func (h *AdminHandler) BanUser(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)

	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID format",
			},
		})
	}

	var req models.BanUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	if len(req.Reason) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Reason must be at least 10 characters",
			},
		})
	}

	ipAddress := c.IP()
	result, err := h.adminUsecase.BanUser(c.Context(), adminID, userID, req.Reason, req.BannedUntil, ipAddress)
	if err != nil {
		log.Printf("Failed to ban user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// UnbanUser handles user unban request
// @Summary Unban user
// @Description Remove active ban from a user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param userId path string true "User ID (UUID)"
// @Param request body models.UnbanUserRequest true "Unban request"
// @Success 200 {object} models.APIResponse{data=models.UnbanUserResponse} "User unbanned"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/users/{userId}/unban [post]
// @Security BearerAuth
func (h *AdminHandler) UnbanUser(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)

	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID format",
			},
		})
	}

	var req models.UnbanUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	if len(req.Reason) < 5 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Reason must be at least 5 characters",
			},
		})
	}

	ipAddress := c.IP()
	result, err := h.adminUsecase.UnbanUser(c.Context(), adminID, userID, req.Reason, ipAddress)
	if err != nil {
		log.Printf("Failed to unban user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetBanHistory handles ban history request
// @Summary Get ban history
// @Description Get ban history for a user (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} models.APIResponse{data=[]models.UserBan} "Ban history"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid user ID"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/users/{userId}/ban-history [get]
// @Security BearerAuth
func (h *AdminHandler) GetBanHistory(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID format",
			},
		})
	}

	bans, err := h.adminUsecase.GetBanHistory(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get ban history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to get ban history",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    bans,
	})
}

// GetAdminActions handles audit log request
// @Summary Get admin actions log
// @Description Get audit log of all admin actions (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param page query int false "Page number" default(1)
// @Param perPage query int false "Results per page (max 100)" default(50)
// @Param action_type query string false "Filter by action type (e.g. CREATE_LEVEL_CONFIG)"
// @Success 200 {object} models.APIResponse{data=models.AdminActionsResponse} "Admin actions"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/actions [get]
// @Security BearerAuth
func (h *AdminHandler) GetAdminActions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.Query("perPage", "50"))
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}

	actionType := c.Query("action_type")

	var (
		result *models.AdminActionsResponse
		err    error
	)
	if actionType != "" {
		result, err = h.adminUsecase.GetActionsFiltered(c.Context(), page, perPage, actionType)
	} else {
		result, err = h.adminUsecase.GetActions(c.Context(), page, perPage)
	}
	if err != nil {
		log.Printf("Failed to get admin actions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to get admin actions",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetSystemStats handles system stats request
// @Summary Get system statistics
// @Description Get system overview statistics (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Success 200 {object} models.APIResponse{data=models.SystemStatsResponse} "System stats"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/stats/overview [get]
// @Security BearerAuth
func (h *AdminHandler) GetSystemStats(c *fiber.Ctx) error {
	stats, err := h.adminUsecase.GetSystemStats(c.Context())
	if err != nil {
		log.Printf("Failed to get system stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to get system statistics",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// ExportUserData handles user data export for GDPR
// @Summary Export user data
// @Description Export all user data for GDPR compliance (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} models.APIResponse{data=models.ExportUserDataResponse} "User data"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid user ID"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/users/{userId}/export-data [get]
// @Security BearerAuth
func (h *AdminHandler) ExportUserData(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID format",
			},
		})
	}

	data, err := h.adminUsecase.ExportUserData(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to export user data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to export user data",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    data,
	})
}

// ResetLeaderboard handles leaderboard reset request
// @Summary Reset leaderboard
// @Description Reset all entries for a specific level (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Param levelId path string true "Level ID"
// @Param request body models.ResetLeaderboardRequest true "Reset reason"
// @Success 200 {object} models.APIResponse "Leaderboard reset"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Invalid request"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Router /admin/leaderboards/{levelId} [delete]
// @Security BearerAuth
func (h *AdminHandler) ResetLeaderboard(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)
	levelID := c.Params("levelId")

	var req models.ResetLeaderboardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	if len(req.Reason) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Reason must be at least 10 characters",
			},
		})
	}

	ipAddress := c.IP()
	err := h.adminUsecase.ResetLeaderboard(c.Context(), adminID, levelID, req.Reason, ipAddress)
	if err != nil {
		log.Printf("Failed to reset leaderboard: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to reset leaderboard",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    fiber.Map{"message": "Leaderboard reset successfully"},
	})
}

// GetLeaderboardStats handles leaderboard stats request
// @Summary Get leaderboard statistics
// @Description Get analytics for level leaderboards (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Success 200 {object} models.APIResponse{data=models.LeaderboardStatsResponse} "Leaderboard stats"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Router /admin/leaderboards/stats [get]
// @Security BearerAuth
func (h *AdminHandler) GetLeaderboardStats(c *fiber.Ctx) error {
	stats, err := h.adminUsecase.GetLeaderboardStats(c.Context())
	if err != nil {
		log.Printf("Failed to get leaderboard stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to get leaderboard statistics",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// @Summary List level configs
// @Tags Admin
// @Produce json
// @Param page query int false "Page" default(1)
// @Param perPage query int false "Per page" default(20)
// @Success 200 {object} models.APIResponse{data=models.LevelConfigListResponse}
// @Router /admin/levels [get]
// @Security BearerAuth
func (h *AdminHandler) ListLevels(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.Query("perPage", "20"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	result, err := h.adminUsecase.ListLevels(c.Context(), page, perPage)
	if err != nil {
		log.Printf("Failed to list levels: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: "Failed to list levels"},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Get level config
// @Tags Admin
// @Produce json
// @Param levelId path string true "Level ID"
// @Success 200 {object} models.APIResponse{data=models.AdminLevelConfig}
// @Router /admin/levels/{levelId} [get]
// @Security BearerAuth
func (h *AdminHandler) GetLevel(c *fiber.Ctx) error {
	levelID := c.Params("levelId")
	result, err := h.adminUsecase.GetLevelConfig(c.Context(), levelID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "NOT_FOUND", Message: err.Error()},
			})
		}
		log.Printf("Failed to get level: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: "Failed to get level"},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Create level config
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body models.CreateLevelConfigRequest true "Level config"
// @Success 201 {object} models.APIResponse{data=models.AdminLevelConfig}
// @Router /admin/levels [post]
// @Security BearerAuth
func (h *AdminHandler) CreateLevel(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)

	var req models.CreateLevelConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Invalid request body"},
		})
	}
	if req.LevelID == "" || req.Name == "" || req.MapID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "level_id, map_id, and name are required"},
		})
	}

	result, err := h.adminUsecase.CreateLevelConfig(c.Context(), adminID, &req, c.IP())
	if err != nil {
		log.Printf("Failed to create level: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: err.Error()},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Update level config
// @Tags Admin
// @Accept json
// @Produce json
// @Param levelId path string true "Level ID"
// @Param request body models.UpdateLevelConfigRequest true "Update fields"
// @Success 200 {object} models.APIResponse{data=models.AdminLevelConfig}
// @Router /admin/levels/{levelId} [put]
// @Security BearerAuth
func (h *AdminHandler) UpdateLevel(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)
	levelID := c.Params("levelId")

	var req models.UpdateLevelConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Invalid request body"},
		})
	}

	result, err := h.adminUsecase.UpdateLevelConfig(c.Context(), adminID, levelID, &req, c.IP())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "NOT_FOUND", Message: err.Error()},
			})
		}
		log.Printf("Failed to update level: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: err.Error()},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Delete level config
// @Tags Admin
// @Produce json
// @Param levelId path string true "Level ID"
// @Success 200 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse{error=models.APIError} "Conflict - has leaderboard entries"
// @Router /admin/levels/{levelId} [delete]
// @Security BearerAuth
func (h *AdminHandler) DeleteLevel(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)
	levelID := c.Params("levelId")

	err := h.adminUsecase.DeleteLevelConfig(c.Context(), adminID, levelID, c.IP())
	if err != nil {
		if strings.HasPrefix(err.Error(), "CONFLICT:") {
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "CONFLICT", Message: err.Error()[len("CONFLICT: "):]},
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "NOT_FOUND", Message: err.Error()},
			})
		}
		log.Printf("Failed to delete level: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: err.Error()},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: fiber.Map{"message": "Level deleted successfully"}})
}

// @Summary List talent configs
// @Tags Admin
// @Produce json
// @Param page query int false "Page" default(1)
// @Param perPage query int false "Per page" default(20)
// @Success 200 {object} models.APIResponse{data=models.TalentConfigListResponse}
// @Router /admin/talents [get]
// @Security BearerAuth
func (h *AdminHandler) ListTalents(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.Query("perPage", "20"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	result, err := h.adminUsecase.ListTalents(c.Context(), page, perPage)
	if err != nil {
		log.Printf("Failed to list talents: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: "Failed to list talents"},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Get talent config
// @Tags Admin
// @Produce json
// @Param talentId path string true "Talent ID"
// @Success 200 {object} models.APIResponse{data=models.AdminTalentConfig}
// @Router /admin/talents/{talentId} [get]
// @Security BearerAuth
func (h *AdminHandler) GetTalent(c *fiber.Ctx) error {
	talentID := c.Params("talentId")
	result, err := h.adminUsecase.GetTalentConfig(c.Context(), talentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "NOT_FOUND", Message: err.Error()},
			})
		}
		log.Printf("Failed to get talent: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: "Failed to get talent"},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Create talent config
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body models.CreateTalentConfigRequest true "Talent config"
// @Success 201 {object} models.APIResponse{data=models.AdminTalentConfig}
// @Router /admin/talents [post]
// @Security BearerAuth
func (h *AdminHandler) CreateTalent(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)

	var req models.CreateTalentConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Invalid request body"},
		})
	}
	if req.TalentID == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "talent_id and name are required"},
		})
	}

	result, err := h.adminUsecase.CreateTalentConfig(c.Context(), adminID, &req, c.IP())
	if err != nil {
		log.Printf("Failed to create talent: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: err.Error()},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Update talent config
// @Tags Admin
// @Accept json
// @Produce json
// @Param talentId path string true "Talent ID"
// @Param request body models.UpdateTalentConfigRequest true "Update fields"
// @Success 200 {object} models.APIResponse{data=models.AdminTalentConfig}
// @Router /admin/talents/{talentId} [put]
// @Security BearerAuth
func (h *AdminHandler) UpdateTalent(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)
	talentID := c.Params("talentId")

	var req models.UpdateTalentConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Invalid request body"},
		})
	}

	result, err := h.adminUsecase.UpdateTalentConfig(c.Context(), adminID, talentID, &req, c.IP())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "NOT_FOUND", Message: err.Error()},
			})
		}
		log.Printf("Failed to update talent: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: err.Error()},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// @Summary Delete talent config
// @Tags Admin
// @Produce json
// @Param talentId path string true "Talent ID"
// @Success 200 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse{error=models.APIError} "Conflict - players have this talent"
// @Router /admin/talents/{talentId} [delete]
// @Security BearerAuth
func (h *AdminHandler) DeleteTalent(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, _ := uuid.Parse(adminIDStr)
	talentID := c.Params("talentId")

	err := h.adminUsecase.DeleteTalentConfig(c.Context(), adminID, talentID, c.IP())
	if err != nil {
		if strings.HasPrefix(err.Error(), "CONFLICT:") {
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "CONFLICT", Message: err.Error()[len("CONFLICT: "):]},
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   &models.APIError{Code: "NOT_FOUND", Message: err.Error()},
			})
		}
		log.Printf("Failed to delete talent: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: err.Error()},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: fiber.Map{"message": "Talent deleted successfully"}})
}

// @Summary Get analytics summary
// @Description Get analytics summary for the last 24h and 7 days (admin only)
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Success 200 {object} models.APIResponse{data=models.AnalyticsSummaryResponse}
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 403 {object} models.APIResponse{error=models.APIError} "Forbidden - Admin only"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /admin/analytics/summary [get]
// @Security BearerAuth
func (h *AdminHandler) GetAnalyticsSummary(c *fiber.Ctx) error {
	result, err := h.adminUsecase.GetAnalyticsSummary(c.Context())
	if err != nil {
		log.Printf("Failed to get analytics summary: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: "Failed to get analytics summary"},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// AdminLogin handles admin dashboard login with username/password
// @Summary Admin login
// @Description Login to admin dashboard with username and password
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body models.AdminLoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.AdminLoginResponse}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 401 {object} models.APIResponse{error=models.APIError}
// @Router /admin/auth/login [post]
func (h *AdminHandler) AdminLogin(c *fiber.Ctx) error {
	var req models.AdminLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Invalid request body"},
		})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Username and password are required"},
		})
	}
	result, err := h.adminUsecase.AdminLogin(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeUnauthorized, Message: err.Error()},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: result})
}

// SetAdminPassword sets password for an admin account
// @Summary Set admin password
// @Description Set or update password for an admin account (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body object true "Password request"
// @Success 200 {object} models.APIResponse
// @Router /admin/auth/set-password [post]
// @Security BearerAuth
func (h *AdminHandler) SetAdminPassword(c *fiber.Ctx) error {
	adminIDStr, _ := c.Locals("userId").(string)
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Invalid user ID"},
		})
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil || len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Password must be at least 8 characters"},
		})
	}

	if err := h.adminUsecase.SetAdminPassword(c.Context(), adminID, req.Password); err != nil {
		log.Printf("Failed to set admin password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   &models.APIError{Code: models.ErrCodeInternalError, Message: "Failed to set password"},
		})
	}
	return c.JSON(models.APIResponse{Success: true, Data: fiber.Map{"message": "Password updated successfully"}})
}
