package repository

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// LevelRepository defines the interface for level-related data access
type LevelRepository interface {
	GetConfig(levelID string, mapID string) (*models.LevelConfig, error)
	GetCompletion(ctx context.Context, userID uuid.UUID, levelID string) (*models.LevelCompletion, error)
	CreateCompletion(ctx context.Context, completion *models.LevelCompletion) error
	UpdateCompletion(ctx context.Context, completion *models.LevelCompletion) error
	UpdateUserStats(ctx context.Context, userID uuid.UUID, goldEarned, starsEarned int) (int, int, error)
}
