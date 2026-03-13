package leaderboard

import (
	"context"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeaderboardUsecase_GetLeaderboard(t *testing.T) {
	repo := new(repo_mock.LeaderboardRepository)
	usecase := NewLeaderboardUsecase(repo)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		entries := []models.HollowWildsLeaderboardEntry{
			{Rank: 1, PlayerID: "P1", Value: 100},
		}
		repo.On("GetByCriteria", ctx, "type", "global", "", 10, 0).Return(entries, 1, nil).Once()

		resp, err := usecase.GetLeaderboard(ctx, "type", "global", "", 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Entries))
		assert.Equal(t, int64(100), resp.Entries[0].Value)
		repo.AssertExpectations(t)
	})
}

func TestLeaderboardUsecase_SubmitEntry(t *testing.T) {
	repo := new(repo_mock.LeaderboardRepository)
	usecase := NewLeaderboardUsecase(repo)

	ctx := context.Background()
	playerID := uuid.New()

	t.Run("new personal best", func(t *testing.T) {
		req := models.LeaderboardSubmitRequest{
			Type:  "score",
			Value: 200,
		}
		repo.On("GetPersonalBest", ctx, playerID, req.Type, "").Return(int64(100), nil).Once()
		repo.On("UpsertEntry", ctx, mock.Anything).Return(nil).Once()

		resp, err := usecase.SubmitEntry(ctx, playerID, req)

		assert.NoError(t, err)
		assert.True(t, resp.IsPersonalBest)
		repo.AssertExpectations(t)
	})

	t.Run("below personal best", func(t *testing.T) {
		req := models.LeaderboardSubmitRequest{
			Type:  "score",
			Value: 50,
		}
		repo.On("GetPersonalBest", ctx, playerID, req.Type, "").Return(int64(100), nil).Once()

		resp, err := usecase.SubmitEntry(ctx, playerID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "value_too_low")
	})
}
