package usecase

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// LeaderboardUsecase defines the interface for leaderboard-related business logic
type LeaderboardUsecase interface {
	// Hollow Wilds (New)
	GetLeaderboard(ctx context.Context, lbType, scope, character string, limit, offset int) (*models.HollowWildsLeaderboardResponse, error)
	SubmitEntry(ctx context.Context, playerID uuid.UUID, req models.LeaderboardSubmitRequest) (*models.LeaderboardSubmitResponse, error)
	GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*models.PlayerLeaderboardResponse, error)

	// Legacy Level Rankings
	GetGlobalLeaderboard(ctx context.Context, levelID string, page, perPage int) (*models.GlobalLeaderboardResponse, error)
	GetPlayerRank(ctx context.Context, userID uuid.UUID, levelID string) (*models.PlayerStatsResponse, error)
	GetFriendsLeaderboard(ctx context.Context, userID uuid.UUID, levelID string) (*models.LevelLeaderboardResponse, error)
	UpdateEntry(ctx context.Context, userID uuid.UUID, levelID string, timeSeconds float64, stars int) error
}
