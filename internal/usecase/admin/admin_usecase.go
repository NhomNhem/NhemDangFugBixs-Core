package admin

import (
	"context"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type adminUsecase struct {
	adminRepo repository.AdminRepository
	lbRepo    repository.LeaderboardRepository
}

// NewAdminUsecase creates a new admin usecase
func NewAdminUsecase(adminRepo repository.AdminRepository, lbRepo repository.LeaderboardRepository) usecase.AdminUsecase {
	return &adminUsecase{
		adminRepo: adminRepo,
		lbRepo:    lbRepo,
	}
}

func (u *adminUsecase) SearchUsers(ctx context.Context, query string, page, perPage int) (*models.UserSearchResponse, error) {
	return u.adminRepo.SearchUsers(ctx, query, page, perPage)
}

func (u *adminUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	return u.adminRepo.GetProfile(ctx, userID)
}

func (u *adminUsecase) AdjustGold(ctx context.Context, adminID, userID uuid.UUID, amount int, reason, ipAddress string) (*models.AdjustGoldResponse, error) {
	return u.adminRepo.AdjustGold(ctx, adminID, userID, amount, reason, ipAddress)
}

func (u *adminUsecase) BanUser(ctx context.Context, adminID, userID uuid.UUID, reason string, bannedUntil *time.Time, ipAddress string) (*models.BanUserResponse, error) {
	return u.adminRepo.BanUser(ctx, adminID, userID, reason, bannedUntil, ipAddress)
}

func (u *adminUsecase) UnbanUser(ctx context.Context, adminID, userID uuid.UUID, reason, ipAddress string) (*models.UnbanUserResponse, error) {
	return u.adminRepo.UnbanUser(ctx, adminID, userID, reason, ipAddress)
}

func (u *adminUsecase) GetBanHistory(ctx context.Context, userID uuid.UUID) ([]models.UserBan, error) {
	return u.adminRepo.GetBanHistory(ctx, userID)
}

func (u *adminUsecase) GetActions(ctx context.Context, page, perPage int) (*models.AdminActionsResponse, error) {
	return u.adminRepo.GetActions(ctx, page, perPage)
}

func (u *adminUsecase) GetSystemStats(ctx context.Context) (*models.SystemStatsResponse, error) {
	return u.adminRepo.GetSystemStats(ctx)
}

func (u *adminUsecase) ExportUserData(ctx context.Context, userID uuid.UUID) (*models.ExportUserDataResponse, error) {
	profile, err := u.adminRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	levelCompletions, _ := u.adminRepo.GetLevelCompletions(ctx, userID)
	talents, _ := u.adminRepo.GetUserTalents(ctx, userID)
	banHistory, _ := u.adminRepo.GetBanHistory(ctx, userID)

	return &models.ExportUserDataResponse{
		User:             *profile,
		LevelCompletions: levelCompletions,
		Talents:          talents,
		BanHistory:       banHistory,
		ExportedAt:       time.Now(),
	}, nil
}

func (u *adminUsecase) ResetLeaderboard(ctx context.Context, adminID uuid.UUID, levelID, reason, ipAddress string) error {
	// 1. Reset in repo
	if err := u.lbRepo.ResetLegacyLeaderboard(ctx, levelID); err != nil {
		return err
	}

	// 2. Log action
	details := map[string]any{
		"level_id": levelID,
		"reason":   reason,
	}
	return u.adminRepo.LogAction(ctx, adminID, "RESET_LEADERBOARD", nil, details, ipAddress)
}

func (u *adminUsecase) GetLeaderboardStats(ctx context.Context) (*models.LeaderboardStatsResponse, error) {
	return u.lbRepo.GetLegacyStats(ctx)
}
