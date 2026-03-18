package usecase

import (
	"context"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// AdminUsecase defines the business logic for administrative operations
type AdminUsecase interface {
	// Auth
	AdminLogin(ctx context.Context, username, password string) (*models.AdminLoginResponse, error)
	SetAdminPassword(ctx context.Context, adminID uuid.UUID, password string) error

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

	// Level Config
	ListLevels(ctx context.Context, page, perPage int) (*models.LevelConfigListResponse, error)
	GetLevelConfig(ctx context.Context, levelID string) (*models.AdminLevelConfig, error)
	CreateLevelConfig(ctx context.Context, adminID uuid.UUID, req *models.CreateLevelConfigRequest, ipAddress string) (*models.AdminLevelConfig, error)
	UpdateLevelConfig(ctx context.Context, adminID uuid.UUID, levelID string, req *models.UpdateLevelConfigRequest, ipAddress string) (*models.AdminLevelConfig, error)
	DeleteLevelConfig(ctx context.Context, adminID uuid.UUID, levelID string, ipAddress string) error

	// Talent Config
	ListTalents(ctx context.Context, page, perPage int) (*models.TalentConfigListResponse, error)
	GetTalentConfig(ctx context.Context, talentID string) (*models.AdminTalentConfig, error)
	CreateTalentConfig(ctx context.Context, adminID uuid.UUID, req *models.CreateTalentConfigRequest, ipAddress string) (*models.AdminTalentConfig, error)
	UpdateTalentConfig(ctx context.Context, adminID uuid.UUID, talentID string, req *models.UpdateTalentConfigRequest, ipAddress string) (*models.AdminTalentConfig, error)
	DeleteTalentConfig(ctx context.Context, adminID uuid.UUID, talentID string, ipAddress string) error

	// Analytics
	GetAnalyticsSummary(ctx context.Context) (*models.AnalyticsSummaryResponse, error)

	// Extended audit log
	GetActionsFiltered(ctx context.Context, page, perPage int, actionType string) (*models.AdminActionsResponse, error)
}
