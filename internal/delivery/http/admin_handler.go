package http

import (
	"log"
	"strconv"

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

	result, err := h.adminUsecase.GetActions(c.Context(), page, perPage)
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
