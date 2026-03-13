package player

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	repository_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlayerUsecase_GetSave(t *testing.T) {
	playerRepo := new(repository_mock.MockPlayerRepository)
	saveRepo := new(repository_mock.MockSaveRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewPlayerUsecase(playerRepo, saveRepo, cacheRepo)

	ctx := context.Background()
	playerID := uuid.New()
	cacheKey := "player:save:" + playerID.String()

	t.Run("cache hit", func(t *testing.T) {
		save := &models.PlayerSave{
			PlayerID: playerID,
			SaveData: models.GameSaveData{
				World: models.WorldData{Seed: 123},
			},
		}
		data, _ := json.Marshal(save)
		cacheRepo.On("Get", ctx, cacheKey).Return(string(data), nil).Once()

		result, err := usecase.GetSave(ctx, playerID)

		assert.NoError(t, err)
		assert.Equal(t, save.SaveData.World.Seed, result.SaveData.World.Seed)
		cacheRepo.AssertExpectations(t)
	})

	t.Run("cache miss - database hit", func(t *testing.T) {
		save := &models.PlayerSave{
			PlayerID: playerID,
			SaveData: models.GameSaveData{
				World: models.WorldData{Seed: 456},
			},
		}
		cacheRepo.On("Get", ctx, cacheKey).Return("", nil).Once()
		saveRepo.On("GetByPlayerID", ctx, playerID).Return(save, nil).Once()
		cacheRepo.On("Set", ctx, cacheKey, mock.Anything, mock.Anything).Return(nil).Once()

		result, err := usecase.GetSave(ctx, playerID)

		assert.NoError(t, err)
		assert.Equal(t, save.SaveData.World.Seed, result.SaveData.World.Seed)
		saveRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
	})
}

func TestPlayerUsecase_SaveGame(t *testing.T) {
	playerRepo := new(repository_mock.MockPlayerRepository)
	saveRepo := new(repository_mock.MockSaveRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewPlayerUsecase(playerRepo, saveRepo, cacheRepo)

	ctx := context.Background()
	playerID := uuid.New()
	saveData := models.GameSaveData{
		World: models.WorldData{Seed: 789},
		Player: models.PlayerState{
			Character: "RIMBA",
			Health:    100,
		},
	}

	t.Run("successful save", func(t *testing.T) {
		saveRepo.On("GetByPlayerID", ctx, playerID).Return(nil, nil).Once()
		saveRepo.On("Upsert", ctx, mock.AnythingOfType("*models.PlayerSave")).Return(nil).Once()
		playerRepo.On("UpdateLastSeen", ctx, playerID).Return(nil).Once()
		cacheRepo.On("Delete", ctx, mock.Anything).Return(nil).Once()

		result, err := usecase.SaveGame(ctx, playerID, saveData, 0)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.SaveVersion)
		saveRepo.AssertExpectations(t)
		playerRepo.AssertExpectations(t)
	})

	t.Run("version conflict", func(t *testing.T) {
		currentSave := &models.PlayerSave{SaveVersion: 5}
		saveRepo.On("GetByPlayerID", ctx, playerID).Return(currentSave, nil).Once()

		_, err := usecase.SaveGame(ctx, playerID, saveData, 3) // Wrong version

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "outdated")
	})
}

func TestPlayerUsecase_RestoreFromBackup(t *testing.T) {
	playerRepo := new(repository_mock.MockPlayerRepository)
	saveRepo := new(repository_mock.MockSaveRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewPlayerUsecase(playerRepo, saveRepo, cacheRepo)

	ctx := context.Background()
	playerID := uuid.New()
	backupID := uuid.New()

	t.Run("successful restore", func(t *testing.T) {
		backup := &models.PlayerSaveBackup{
			ID:       backupID,
			PlayerID: playerID,
			SaveData: models.GameSaveData{World: models.WorldData{Seed: 1}},
		}
		saveRepo.On("GetBackupByID", ctx, backupID).Return(backup, nil).Once()
		saveRepo.On("GetByPlayerID", ctx, playerID).Return(&models.PlayerSave{SaveVersion: 10}, nil).Once()
		saveRepo.On("Upsert", ctx, mock.MatchedBy(func(s *models.PlayerSave) bool {
			return s.SaveVersion == 11
		})).Return(nil).Once()
		cacheRepo.On("Delete", ctx, mock.Anything).Return(nil).Once()

		result, err := usecase.RestoreFromBackup(ctx, playerID, backupID)

		assert.NoError(t, err)
		assert.Equal(t, 11, result.SaveVersion)
	})
}
