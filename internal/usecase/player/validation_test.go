package player

import (
	"context"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlayerUsecase_Validation(t *testing.T) {
	saveRepo := new(repo_mock.SaveRepository)
	cacheRepo := new(repo_mock.CacheRepository)
	usecase := NewPlayerUsecase(saveRepo, cacheRepo)

	ctx := context.Background()
	playerID := uuid.New()

	t.Run("invalid health - too high", func(t *testing.T) {
		invalidData := models.GameSaveData{
			Player: models.PlayerState{
				Character: "RIMBA",
				Health:    150, // Max is 100
			},
		}

		result, err := usecase.SaveGame(ctx, playerID, invalidData, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("invalid character", func(t *testing.T) {
		invalidData := models.GameSaveData{
			Player: models.PlayerState{
				Character: "HACKER",
				Health:    100,
			},
		}

		result, err := usecase.SaveGame(ctx, playerID, invalidData, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("valid data passes", func(t *testing.T) {
		validData := models.GameSaveData{
			Player: models.PlayerState{
				Character: "RIMBA",
				Health:    80,
				Hunger:    50,
				Sanity:    100,
				Warmth:    100,
			},
		}

		saveRepo.On("GetByPlayerID", ctx, playerID).Return(nil, nil).Once()
		saveRepo.On("Upsert", ctx, mock.Anything).Return(nil).Once()
		cacheRepo.On("Delete", ctx, mock.Anything).Return(nil).Once()

		result, err := usecase.SaveGame(ctx, playerID, validData, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		saveRepo.AssertExpectations(t)
	})
}
