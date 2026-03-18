package admin

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type adminUsecase struct {
	adminRepo repository.AdminRepository
	lbRepo    repository.LeaderboardRepository
}

func NewAdminUsecase(adminRepo repository.AdminRepository, lbRepo repository.LeaderboardRepository) usecase.AdminUsecase {
	return &adminUsecase{adminRepo: adminRepo, lbRepo: lbRepo}
}

func (u *adminUsecase) AdminLogin(ctx context.Context, username, password string) (*models.AdminLoginResponse, error) {
	profile, passwordHash, err := u.adminRepo.GetAdminByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if passwordHash == "" {
		return nil, fmt.Errorf("password login not configured for this account")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := u.generateAdminJWT(profile.ID.String(), profile.PlayFabID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AdminLoginResponse{JWT: token, Admin: *profile}, nil
}

func (u *adminUsecase) SetAdminPassword(ctx context.Context, adminID uuid.UUID, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return u.adminRepo.SetAdminPassword(ctx, adminID, string(hash))
}

func (u *adminUsecase) generateAdminJWT(userID, playfabID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-123"
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"userId":    userID,
		"sub":       userID,
		"playfabId": playfabID,
		"iat":       now.Unix(),
		"exp":       now.Add(8 * time.Hour).Unix(),
		"jti":       uuid.New().String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
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
	if err := u.lbRepo.ResetLegacyLeaderboard(ctx, levelID); err != nil {
		return err
	}
	details := map[string]any{"level_id": levelID, "reason": reason}
	return u.adminRepo.LogAction(ctx, adminID, "RESET_LEADERBOARD", nil, details, ipAddress)
}

func (u *adminUsecase) GetLeaderboardStats(ctx context.Context) (*models.LeaderboardStatsResponse, error) {
	return u.lbRepo.GetLegacyStats(ctx)
}

func (u *adminUsecase) ListLevels(ctx context.Context, page, perPage int) (*models.LevelConfigListResponse, error) {
	return u.adminRepo.ListLevelConfigs(ctx, page, perPage)
}

func (u *adminUsecase) GetLevelConfig(ctx context.Context, levelID string) (*models.AdminLevelConfig, error) {
	return u.adminRepo.GetLevelConfig(ctx, levelID)
}

func (u *adminUsecase) CreateLevelConfig(ctx context.Context, adminID uuid.UUID, req *models.CreateLevelConfigRequest, ipAddress string) (*models.AdminLevelConfig, error) {
	config := &models.AdminLevelConfig{
		LevelID: req.LevelID, MapID: req.MapID, Name: req.Name,
		Difficulty: req.Difficulty, MinTimeSeconds: req.MinTimeSeconds,
		BaseGold: req.BaseGold, RewardStars: req.RewardStars,
		Objectives: req.Objectives, IsActive: req.IsActive,
	}
	if err := u.adminRepo.CreateLevelConfig(ctx, config); err != nil {
		return nil, fmt.Errorf("create level config: %w", err)
	}
	u.adminRepo.LogAction(ctx, adminID, "CREATE_LEVEL_CONFIG", nil, map[string]any{"level_id": req.LevelID, "name": req.Name}, ipAddress)
	return u.adminRepo.GetLevelConfig(ctx, req.LevelID)
}

func (u *adminUsecase) UpdateLevelConfig(ctx context.Context, adminID uuid.UUID, levelID string, req *models.UpdateLevelConfigRequest, ipAddress string) (*models.AdminLevelConfig, error) {
	existing, err := u.adminRepo.GetLevelConfig(ctx, levelID)
	if err != nil {
		return nil, fmt.Errorf("level config not found: %w", err)
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Difficulty != nil {
		existing.Difficulty = *req.Difficulty
	}
	if req.MinTimeSeconds != nil {
		existing.MinTimeSeconds = *req.MinTimeSeconds
	}
	if req.BaseGold != nil {
		existing.BaseGold = *req.BaseGold
	}
	if req.RewardStars != nil {
		existing.RewardStars = *req.RewardStars
	}
	if req.Objectives != nil {
		existing.Objectives = req.Objectives
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if err := u.adminRepo.UpdateLevelConfig(ctx, existing); err != nil {
		return nil, fmt.Errorf("update level config: %w", err)
	}
	u.adminRepo.LogAction(ctx, adminID, "UPDATE_LEVEL_CONFIG", nil, map[string]any{"level_id": levelID}, ipAddress)
	return u.adminRepo.GetLevelConfig(ctx, levelID)
}

func (u *adminUsecase) DeleteLevelConfig(ctx context.Context, adminID uuid.UUID, levelID string, ipAddress string) error {
	count, err := u.adminRepo.LevelHasLeaderboardEntries(ctx, levelID)
	if err != nil {
		return fmt.Errorf("check leaderboard entries: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("CONFLICT: Cannot delete level: %d leaderboard entries exist. Reset the leaderboard first.", count)
	}
	if err := u.adminRepo.DeleteLevelConfig(ctx, levelID); err != nil {
		return fmt.Errorf("delete level config: %w", err)
	}
	u.adminRepo.LogAction(ctx, adminID, "DELETE_LEVEL_CONFIG", nil, map[string]any{"level_id": levelID}, ipAddress)
	return nil
}

func (u *adminUsecase) ListTalents(ctx context.Context, page, perPage int) (*models.TalentConfigListResponse, error) {
	return u.adminRepo.ListTalentConfigs(ctx, page, perPage)
}

func (u *adminUsecase) GetTalentConfig(ctx context.Context, talentID string) (*models.AdminTalentConfig, error) {
	return u.adminRepo.GetTalentConfig(ctx, talentID)
}

func (u *adminUsecase) CreateTalentConfig(ctx context.Context, adminID uuid.UUID, req *models.CreateTalentConfigRequest, ipAddress string) (*models.AdminTalentConfig, error) {
	config := &models.AdminTalentConfig{
		TalentID: req.TalentID, Name: req.Name, Description: req.Description,
		MaxLevel: req.MaxLevel, BaseCost: req.BaseCost, CostScaling: req.CostScaling,
		BonusPerLevel: req.BonusPerLevel, StatType: req.StatType,
		UnlockMap: req.UnlockMap, IsActive: req.IsActive,
	}
	if err := u.adminRepo.CreateTalentConfig(ctx, config); err != nil {
		return nil, fmt.Errorf("create talent config: %w", err)
	}
	u.adminRepo.LogAction(ctx, adminID, "CREATE_TALENT_CONFIG", nil, map[string]any{"talent_id": req.TalentID, "name": req.Name}, ipAddress)
	return u.adminRepo.GetTalentConfig(ctx, req.TalentID)
}

func (u *adminUsecase) UpdateTalentConfig(ctx context.Context, adminID uuid.UUID, talentID string, req *models.UpdateTalentConfigRequest, ipAddress string) (*models.AdminTalentConfig, error) {
	existing, err := u.adminRepo.GetTalentConfig(ctx, talentID)
	if err != nil {
		return nil, fmt.Errorf("talent config not found: %w", err)
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.MaxLevel != nil {
		existing.MaxLevel = *req.MaxLevel
	}
	if req.BaseCost != nil {
		existing.BaseCost = *req.BaseCost
	}
	if req.CostScaling != nil {
		existing.CostScaling = *req.CostScaling
	}
	if req.BonusPerLevel != nil {
		existing.BonusPerLevel = *req.BonusPerLevel
	}
	if req.StatType != nil {
		existing.StatType = *req.StatType
	}
	if req.UnlockMap != nil {
		existing.UnlockMap = *req.UnlockMap
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if err := u.adminRepo.UpdateTalentConfig(ctx, existing); err != nil {
		return nil, fmt.Errorf("update talent config: %w", err)
	}
	u.adminRepo.LogAction(ctx, adminID, "UPDATE_TALENT_CONFIG", nil, map[string]any{"talent_id": talentID}, ipAddress)
	return u.adminRepo.GetTalentConfig(ctx, talentID)
}

func (u *adminUsecase) DeleteTalentConfig(ctx context.Context, adminID uuid.UUID, talentID string, ipAddress string) error {
	count, err := u.adminRepo.TalentHasPlayers(ctx, talentID)
	if err != nil {
		return fmt.Errorf("check talent players: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("CONFLICT: Cannot delete talent: %d players have this talent unlocked.", count)
	}
	if err := u.adminRepo.DeleteTalentConfig(ctx, talentID); err != nil {
		return fmt.Errorf("delete talent config: %w", err)
	}
	u.adminRepo.LogAction(ctx, adminID, "DELETE_TALENT_CONFIG", nil, map[string]any{"talent_id": talentID}, ipAddress)
	return nil
}

func (u *adminUsecase) GetAnalyticsSummary(ctx context.Context) (*models.AnalyticsSummaryResponse, error) {
	return u.adminRepo.GetAnalyticsSummary(ctx)
}

func (u *adminUsecase) GetActionsFiltered(ctx context.Context, page, perPage int, actionType string) (*models.AdminActionsResponse, error) {
	return u.adminRepo.GetActionsFiltered(ctx, page, perPage, actionType)
}
