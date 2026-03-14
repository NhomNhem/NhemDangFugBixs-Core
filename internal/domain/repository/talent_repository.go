package repository

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// TalentRepository defines the interface for talent-related data access
type TalentRepository interface {
	GetConfigs(ctx context.Context) (map[string]*models.TalentConfig, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserTalent, error)
	GetOrCreate(ctx context.Context, userID uuid.UUID, talentID string) (*models.UserTalent, error)
	UpdateLevel(ctx context.Context, id uuid.UUID, newLevel int) error
	UpdateUserGold(ctx context.Context, userID uuid.UUID, goldChange int) (int, error)
}
