package repository

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// SaveRepository defines the interface for player save data access
type SaveRepository interface {
	GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*models.PlayerSave, error)
	Upsert(ctx context.Context, save *models.PlayerSave) error
	CreateBackup(ctx context.Context, backup *models.PlayerSaveBackup) error
	GetBackupsByPlayerID(ctx context.Context, playerID uuid.UUID) ([]models.PlayerSaveBackup, error)
	CountBackups(ctx context.Context, playerID uuid.UUID) (int, error)
	DeleteOldestBackup(ctx context.Context, playerID uuid.UUID) error
}
