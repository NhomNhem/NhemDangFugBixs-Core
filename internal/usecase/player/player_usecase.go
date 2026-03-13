package player

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type playerUsecase struct {
	saveRepo  repository.SaveRepository
	cacheRepo repository.CacheRepository
}

// NewPlayerUsecase creates a new player usecase
func NewPlayerUsecase(
	saveRepo repository.SaveRepository,
	cacheRepo repository.CacheRepository,
) usecase.PlayerUsecase {
	return &playerUsecase{
		saveRepo:  saveRepo,
		cacheRepo: cacheRepo,
	}
}

func (u *playerUsecase) GetSave(ctx context.Context, playerID uuid.UUID) (*models.PlayerSave, error) {
	// 1. Try cache
	cacheKey := fmt.Sprintf("player:save:%s", playerID.String())
	cached, err := u.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var save models.PlayerSave
		if err := json.Unmarshal([]byte(cached), &save); err == nil {
			return &save, nil
		}
	}

	// 2. Query database
	save, err := u.saveRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	if save != nil {
		// 3. Cache the result
		if saveDataStr, err := json.Marshal(save); err == nil {
			u.cacheRepo.Set(ctx, cacheKey, string(saveDataStr), 5*time.Minute)
		}
	}

	return save, nil
}

func (u *playerUsecase) SaveGame(ctx context.Context, playerID uuid.UUID, saveData models.GameSaveData, expectedVersion int) (*models.PlayerSave, error) {
	// 1. Get current version for conflict check
	currentSave, err := u.saveRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	currentVersion := 0
	if currentSave != nil {
		currentVersion = currentSave.SaveVersion
	}

	// 2. Optimistic locking check
	if currentSave != nil && expectedVersion != 0 && expectedVersion != currentVersion {
		return nil, &models.VersionConflictError{
			ErrorCode:     "version_conflict",
			ServerVersion: currentVersion,
			Message:       "Save is outdated, fetch latest first",
		}
	}

	// 3. Prepare new save
	newSave := &models.PlayerSave{
		PlayerID:    playerID,
		SaveVersion: currentVersion + 1,
		SaveData:    saveData,
	}

	// 4. Persist
	if err := u.saveRepo.Upsert(ctx, newSave); err != nil {
		return nil, err
	}

	// 5. Handle automatic backup (every 10 versions)
	if newSave.SaveVersion%10 == 0 {
		u.CreateBackup(ctx, playerID)
	}

	// 6. Invalidate cache
	cacheKey := fmt.Sprintf("player:save:%s", playerID.String())
	u.cacheRepo.Delete(ctx, cacheKey)

	return newSave, nil
}

func (u *playerUsecase) CreateBackup(ctx context.Context, playerID uuid.UUID) (*models.PlayerSaveBackup, error) {
	// 1. Get latest save
	save, err := u.saveRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if save == nil {
		return nil, fmt.Errorf("no save data found to backup")
	}

	// 2. Check backup count limit (10)
	count, err := u.saveRepo.CountBackups(ctx, playerID)
	if err == nil && count >= 10 {
		u.saveRepo.DeleteOldestBackup(ctx, playerID)
	}

	// 3. Create backup
	backup := &models.PlayerSaveBackup{
		PlayerID:    playerID,
		SaveVersion: save.SaveVersion,
		SaveData:    save.SaveData,
	}

	if err := u.saveRepo.CreateBackup(ctx, backup); err != nil {
		return nil, err
	}

	return backup, nil
}

func (u *playerUsecase) GetBackups(ctx context.Context, playerID uuid.UUID) ([]models.PlayerSaveBackup, error) {
	return u.saveRepo.GetBackupsByPlayerID(ctx, playerID)
}
