package usecase

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// TalentUsecase defines the business logic for character talents
type TalentUsecase interface {
	GetTalentConfigs() map[string]*models.TalentConfig
	GetUserTalents(ctx context.Context, userID uuid.UUID) ([]models.UserTalent, error)
	UpgradeTalent(ctx context.Context, userID uuid.UUID, talentID string) (*models.TalentUpgradeResponse, error)
}
