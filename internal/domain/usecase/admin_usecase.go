package usecase

import (
	"context"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// AdminUsecase defines the business logic for administrative operations
type AdminUsecase interface {
	SearchUsers(ctx context.Context, query string, page, perPage int) (*models.UserSearchResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	AdjustGold(ctx context.Context, adminID, userID uuid.UUID, amount int, reason, ipAddress string) (*models.AdjustGoldResponse, error)
	BanUser(ctx context.Context, adminID, userID uuid.UUID, reason string, bannedUntil *time.Time, ipAddress string) (*models.BanUserResponse, error)
	UnbanUser(ctx context.Context, adminID, userID uuid.UUID, reason, ipAddress string) (*models.UnbanUserResponse, error)
	GetBanHistory(ctx context.Context, userID uuid.UUID) ([]models.UserBan, error)
	GetActions(ctx context.Context, page, perPage int) (*models.AdminActionsResponse, error)
	GetSystemStats(ctx context.Context) (*models.SystemStatsResponse, error)
	ExportUserData(ctx context.Context, userID uuid.UUID) (*models.ExportUserDataResponse, error)

	// Leaderboard Management
	ResetLeaderboard(ctx context.Context, adminID uuid.UUID, levelID, reason, ipAddress string) error
	GetLeaderboardStats(ctx context.Context) (*models.LeaderboardStatsResponse, error)
}
