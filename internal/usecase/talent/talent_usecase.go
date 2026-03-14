package talent

import (
	"context"
	"fmt"
	"math"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type talentUsecase struct {
	talentRepo repository.TalentRepository
}

// NewTalentUsecase creates a new talent usecase
func NewTalentUsecase(talentRepo repository.TalentRepository) usecase.TalentUsecase {
	return &talentUsecase{
		talentRepo: talentRepo,
	}
}

func (u *talentUsecase) GetTalentConfigs(ctx context.Context) (map[string]*models.TalentConfig, error) {
	return u.talentRepo.GetConfigs(ctx)
}

func (u *talentUsecase) GetUserTalents(ctx context.Context, userID uuid.UUID) ([]models.UserTalent, error) {
	return u.talentRepo.GetByUserID(ctx, userID)
}

func (u *talentUsecase) UpgradeTalent(ctx context.Context, userID uuid.UUID, talentID string) (*models.TalentUpgradeResponse, error) {
	// 1. Get talent config
	configs, err := u.talentRepo.GetConfigs(ctx)
	if err != nil {
		return nil, err
	}
	config, ok := configs[talentID]
	if !ok {
		return nil, fmt.Errorf("invalid talent ID: %s", talentID)
	}

	// 2. Get or create user talent
	talent, err := u.talentRepo.GetOrCreate(ctx, userID, talentID)
	if err != nil {
		return nil, err
	}

	// 3. Check if already max level
	if talent.CurrentLevel >= config.MaxLevel {
		return nil, fmt.Errorf("talent already at max level")
	}

	// 4. Calculate upgrade cost
	// Formula: baseCost * (scaling ^ currentLevel)
	cost := float64(config.BaseCost) * math.Pow(config.CostScaling, float64(talent.CurrentLevel))
	upgradeCost := int(math.Round(cost))

	// 5. Update user gold (deduct)
	newTotalGold, err := u.talentRepo.UpdateUserGold(ctx, userID, -upgradeCost)
	if err != nil {
		return nil, err
	}

	// 6. Upgrade talent level
	newLevel := talent.CurrentLevel + 1
	if err := u.talentRepo.UpdateLevel(ctx, talent.ID, newLevel); err != nil {
		// NOTE: In a real production system, we'd want more robust transaction handling here.
		// For now, following the established pattern.
		return nil, err
	}

	// 7. Prepare response
	totalBonus := config.BonusPerLevel * float64(newLevel)
	statBonus := fmt.Sprintf("+%.0f%% %s", totalBonus, config.StatType)

	var nextLevelCost int
	if newLevel < config.MaxLevel {
		nextCost := float64(config.BaseCost) * math.Pow(config.CostScaling, float64(newLevel))
		nextLevelCost = int(math.Round(nextCost))
	}

	return &models.TalentUpgradeResponse{
		Success:       true,
		TalentID:      talentID,
		NewLevel:      newLevel,
		GoldSpent:     upgradeCost,
		NewTotalGold:  newTotalGold,
		StatBonus:     statBonus,
		NextLevelCost: nextLevelCost,
	}, nil
}
