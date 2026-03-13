package usecase

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// LeaderboardUsecase defines the interface for leaderboard-related business logic
type LeaderboardUsecase interface {
	GetLeaderboard(ctx context.Context, lbType, scope, character string, limit, offset int) (*models.HollowWildsLeaderboardResponse, error)
	SubmitEntry(ctx context.Context, playerID uuid.UUID, req models.LeaderboardSubmitRequest) (*models.LeaderboardSubmitResponse, error)
	GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*models.PlayerLeaderboardResponse, error)
}
