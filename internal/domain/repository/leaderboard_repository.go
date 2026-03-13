package repository

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// LeaderboardRepository defines the interface for leaderboard data access
type LeaderboardRepository interface {
	GetByCriteria(ctx context.Context, lbType, scope, character string, limit, offset int) ([]models.HollowWildsLeaderboardEntry, int, error)
	GetPersonalBest(ctx context.Context, playerID uuid.UUID, lbType, character string) (int64, error)
	UpsertEntry(ctx context.Context, entry *models.HollowWildsLeaderboardEntry) error
	GetPlayerStats(ctx context.Context, playerID uuid.UUID) ([]models.PlayerLeaderboardEntry, error)
}
