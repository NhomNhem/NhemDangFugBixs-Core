package http

import (
	"log"
	"strings"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LeaderboardHandler handles leaderboard endpoints
type LeaderboardHandler struct {
	leaderboardUsecase usecase.LeaderboardUsecase
}

// NewLeaderboardHandler creates a new leaderboard handler
func NewLeaderboardHandler(leaderboardUsecase usecase.LeaderboardUsecase) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardUsecase: leaderboardUsecase,
	}
}

// GetHollowWildsLeaderboard handles the new Hollow Wilds leaderboard request
// @Summary Get HW Leaderboard
// @Description Get ranked entries for a specific metric and scope
// @Tags HollowWilds
// @Produce json
// @Param type query string false "Metric type (longest_run_days, sebilah_soul_level, bosses_killed)" default(longest_run_days)
// @Param scope query string false "Scope (global, per_character)" default(global)
// @Param character query string false "Character filter (required if scope=per_character)"
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.HollowWildsLeaderboardResponse "Leaderboard data"
// @Router /leaderboard [get]
func (h *LeaderboardHandler) GetHollowWildsLeaderboard(c *fiber.Ctx) error {
	lbType := c.Query("type", "longest_run_days")
	scope := c.Query("scope", "global")
	character := c.Query("character", "")
	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	if scope == "per_character" && character == "" {
		return fiber.NewError(fiber.StatusBadRequest, "character is required for per_character scope")
	}

	leaderboard, err := h.leaderboardUsecase.GetLeaderboard(c.Context(), lbType, scope, character, limit, offset)
	if err != nil {
		log.Printf("Failed to get Hollow Wilds leaderboard: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve leaderboard")
	}

	return c.JSON(leaderboard)
}

// SubmitHollowWildsEntry handles leaderboard submission
// @Summary Submit HW Run
// @Description Submit a result after a run to update personal best and ranks
// @Tags HollowWilds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.LeaderboardSubmitRequest true "Run result"
// @Success 200 {object} models.LeaderboardSubmitResponse "Submission result"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Value too low or invalid request"
// @Router /leaderboard/submit [post]
func (h *LeaderboardHandler) SubmitHollowWildsEntry(c *fiber.Ctx) error {
	playerIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return fiber.ErrUnauthorized
	}
	playerID, _ := uuid.Parse(playerIDStr)

	var req models.LeaderboardSubmitRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.leaderboardUsecase.SubmitEntry(c.Context(), playerID, req)
	if err != nil {
		log.Printf("Failed to submit leaderboard entry: %v", err)
		if strings.Contains(err.Error(), "value_too_low") {
			return fiber.NewError(fiber.StatusBadRequest, "Submitted value does not beat personal best")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to submit entry")
	}

	return c.JSON(result)
}

// GetPlayerHollowWildsStats handles request for player's own ranks
// @Summary Get HW Player Ranks
// @Description Get current ranks for the authenticated player across all types
// @Tags HollowWilds
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.PlayerLeaderboardResponse "Player rankings"
// @Router /leaderboard/player [get]
func (h *LeaderboardHandler) GetPlayerHollowWildsStats(c *fiber.Ctx) error {
	playerIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return fiber.ErrUnauthorized
	}
	playerID, _ := uuid.Parse(playerIDStr)

	stats, err := h.leaderboardUsecase.GetPlayerStats(c.Context(), playerID)
	if err != nil {
		log.Printf("Failed to get player Hollow Wilds stats: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve player stats")
	}

	return c.JSON(stats)
}

// GetGlobalLeaderboard handles the request for level-specific global rankings
// @Summary (DEPRECATED) Get Global Level Leaderboard
// @Description [DEPRECATED] Get ranked times for a specific level. Use /api/v1/leaderboard instead.
// @Tags Leaderboard (Legacy)
// @Param levelId path string true "Level ID"
// @Param page query int false "Page number" default(1)
// @Param perPage query int false "Entries per page" default(10)
// @Success 200 {object} models.GlobalLeaderboardResponse
// @Router /leaderboards/{levelId} [get]
func (h *LeaderboardHandler) GetGlobalLeaderboard(c *fiber.Ctx) error {
	log.Printf("[DEPRECATION WARNING] Legacy GetGlobalLeaderboard called for level: %s", c.Params("levelId"))
	levelID := c.Params("levelId")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("perPage", 10)

	resp, err := h.leaderboardUsecase.GetGlobalLeaderboard(c.Context(), levelID, page, perPage)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve global leaderboard")
	}

	return c.JSON(resp)
}

// GetPlayerRank handles the request for player's rank on a specific level
// @Summary (DEPRECATED) Get Player Level Rank
// @Description [DEPRECATED] Get player's rank and surrounding players for a specific level.
// @Tags Leaderboard (Legacy)
// @Security BearerAuth
// @Param levelId path string true "Level ID"
// @Success 200 {object} models.PlayerStatsResponse
// @Router /leaderboards/{levelId}/me [get]
func (h *LeaderboardHandler) GetPlayerRank(c *fiber.Ctx) error {
	log.Printf("[DEPRECATION WARNING] Legacy GetPlayerRank called for level: %s", c.Params("levelId"))
	userIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return fiber.ErrUnauthorized
	}
	userID, _ := uuid.Parse(userIDStr)
	levelID := c.Params("levelId")

	resp, err := h.leaderboardUsecase.GetPlayerRank(c.Context(), userID, levelID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve player rank")
	}

	return c.JSON(resp)
}

// GetFriendsLeaderboard handles the request for friends rankings on a specific level
// @Summary (DEPRECATED) Get Friends Level Leaderboard
// @Description [DEPRECATED] Get friends rankings for a specific level.
// @Tags Leaderboard (Legacy)
// @Security BearerAuth
// @Param levelId path string true "Level ID"
// @Success 200 {object} models.LevelLeaderboardResponse
// @Router /leaderboards/{levelId}/friends [get]
func (h *LeaderboardHandler) GetFriendsLeaderboard(c *fiber.Ctx) error {
	log.Printf("[DEPRECATION WARNING] Legacy GetFriendsLeaderboard called for level: %s", c.Params("levelId"))
	userIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return fiber.ErrUnauthorized
	}
	userID, _ := uuid.Parse(userIDStr)
	levelID := c.Params("levelId")

	resp, err := h.leaderboardUsecase.GetFriendsLeaderboard(c.Context(), userID, levelID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve friends leaderboard")
	}

	return c.JSON(resp)
}
