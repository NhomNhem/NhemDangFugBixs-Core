package usecase

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// PlayerUsecase defines the interface for player-related business logic (Save/Load/Backup)
type PlayerUsecase interface {
	GetSave(ctx context.Context, playerID uuid.UUID) (*models.PlayerSave, error)
	SaveGame(ctx context.Context, playerID uuid.UUID, saveData models.GameSaveData, expectedVersion int) (*models.PlayerSave, error)
	CreateBackup(ctx context.Context, playerID uuid.UUID) (*models.PlayerSaveBackup, error)
	GetBackups(ctx context.Context, playerID uuid.UUID) ([]models.PlayerSaveBackup, error)
}
