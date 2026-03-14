package level

import (
	"context"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	repository_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	usecase_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLevelUsecase_CompleteLevel(t *testing.T) {
	levelRepo := new(repository_mock.MockLevelRepository)
	lbUsecase := new(usecase_mock.MockLeaderboardUsecase)
	usecase := NewLevelUsecase(levelRepo, lbUsecase)

	ctx := context.Background()
	userID := uuid.New()
	req := &models.LevelCompletionRequest{
		LevelID:         "level_1",
		MapID:           "map_1",
		TimeSeconds:     45.0,
		FinalHP:         80.0,
		DashCount:       5,
		CounterCount:    2,
		VulnerableKills: 3,
	}

	config := &models.LevelConfig{
		LevelID:        "level_1",
		MapID:          "map_1",
		MinTimeSeconds: 10.0,
		BaseGold:       100,
		Objectives: []models.LevelObjective{
			{Type: "completion", Threshold: 1, Operator: "gte"},
			{Type: "health", Threshold: 50, Operator: "gte"},
			{Type: "time", Threshold: 60, Operator: "lte"},
		},
	}

	t.Run("successful first completion", func(t *testing.T) {
		levelRepo.On("GetConfig", ctx, req.LevelID, req.MapID).Return(config, nil).Once()
		levelRepo.On("GetCompletion", ctx, userID, req.LevelID).Return(nil, nil).Once()
		levelRepo.On("CreateCompletion", ctx, mock.AnythingOfType("*models.LevelCompletion")).Return(nil).Once()

		// stars=3, gold=100 + 3*20 = 160
		levelRepo.On("UpdateUserStats", ctx, userID, 160, 3).Return(1160, 3, nil).Once()

		// Leaderboard update
		lbUsecase.On("UpdateEntry", ctx, userID, req.LevelID, req.TimeSeconds, 3).Return(nil).Once()

		resp, err := usecase.CompleteLevel(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
		assert.Equal(t, 3, resp.StarsEarned)
		assert.Equal(t, 160, resp.GoldEarned)
		assert.True(t, resp.IsFirstCompletion)

		levelRepo.AssertExpectations(t)
		lbUsecase.AssertExpectations(t)
	})

	t.Run("anti-cheat severe - too fast", func(t *testing.T) {
		reqTooFast := *req
		reqTooFast.TimeSeconds = 5.0 // Below config.MinTimeSeconds (10.0)

		levelRepo.On("GetConfig", ctx, req.LevelID, req.MapID).Return(config, nil).Once()

		resp, err := usecase.CompleteLevel(ctx, userID, &reqTooFast)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid completion data", err.Error())
	})
}
