package repository

import (
	"context"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// AdminRepository defines the interface for administrative data access
type AdminRepository interface {
	SearchUsers(ctx context.Context, query string, page, perPage int) (*models.UserSearchResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	AdjustGold(ctx context.Context, adminID, userID uuid.UUID, amount int, reason, ipAddress string) (*models.AdjustGoldResponse, error)
	BanUser(ctx context.Context, adminID, userID uuid.UUID, reason string, bannedUntil *time.Time, ipAddress string) (*models.BanUserResponse, error)
	UnbanUser(ctx context.Context, adminID, userID uuid.UUID, reason, ipAddress string) (*models.UnbanUserResponse, error)
	GetBanHistory(ctx context.Context, userID uuid.UUID) ([]models.UserBan, error)
	GetActions(ctx context.Context, page, perPage int) (*models.AdminActionsResponse, error)
	GetSystemStats(ctx context.Context) (*models.SystemStatsResponse, error)
	GetLevelCompletions(ctx context.Context, userID uuid.UUID) ([]map[string]any, error)
	GetUserTalents(ctx context.Context, userID uuid.UUID) ([]map[string]any, error)
	LogAction(ctx context.Context, adminID uuid.UUID, actionType string, targetUserID *uuid.UUID, details map[string]any, ipAddress string) error
}
