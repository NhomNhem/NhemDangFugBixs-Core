package player

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlayerUsecase_GetSave(t *testing.T) {
	saveRepo := new(repo_mock.SaveRepository)
	cacheRepo := new(repo_mock.CacheRepository)
	usecase := NewPlayerUsecase(saveRepo, cacheRepo)

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
