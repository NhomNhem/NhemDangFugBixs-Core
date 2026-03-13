package repository

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// LeaderboardRepository defines the interface for leaderboard data access
type LeaderboardRepository interface {
	// Hollow Wilds (New)
	GetByCriteria(ctx context.Context, lbType, scope, character string, limit, offset int) ([]models.HollowWildsLeaderboardEntry, int, error)
	GetPersonalBest(ctx context.Context, playerID uuid.UUID, lbType, character string) (int64, error)
	UpsertEntry(ctx context.Context, entry *models.HollowWildsLeaderboardEntry) error
	GetPlayerStats(ctx context.Context, playerID uuid.UUID) ([]models.PlayerLeaderboardEntry, error)

	// Legacy Level Rankings
	GetLegacyGlobal(ctx context.Context, levelID string, limit, offset int) ([]models.LeaderboardEntry, int, error)
	GetLegacyPlayerRank(ctx context.Context, userID uuid.UUID, levelID string) (int, float64, int, error)
	UpsertLegacyEntry(ctx context.Context, userID uuid.UUID, levelID string, timeSeconds float64, stars int) error
	GetLegacyFriends(ctx context.Context, friendIDs []string, levelID string) ([]models.LeaderboardEntry, error)

	// Admin
	ResetLegacyLeaderboard(ctx context.Context, levelID string) error
	GetLegacyStats(ctx context.Context) (*models.LeaderboardStatsResponse, error)
}
