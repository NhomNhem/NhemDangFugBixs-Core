package leaderboard

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	repository_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeaderboardUsecase_GetLeaderboard(t *testing.T) {
	repo := new(repository_mock.MockLeaderboardRepository)
	playerRepo := new(repository_mock.MockPlayerRepository)
	identityRepo := new(repository_mock.MockIdentityRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewLeaderboardUsecase(repo, playerRepo, identityRepo, cacheRepo)

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
	repo := new(repository_mock.MockLeaderboardRepository)
	playerRepo := new(repository_mock.MockPlayerRepository)
	identityRepo := new(repository_mock.MockIdentityRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewLeaderboardUsecase(repo, playerRepo, identityRepo, cacheRepo)

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
}

func TestLeaderboardUsecase_GetGlobalLeaderboard(t *testing.T) {
	repo := new(repository_mock.MockLeaderboardRepository)
	playerRepo := new(repository_mock.MockPlayerRepository)
	identityRepo := new(repository_mock.MockIdentityRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewLeaderboardUsecase(repo, playerRepo, identityRepo, cacheRepo)

	ctx := context.Background()
	levelID := "level1"

	t.Run("successful retrieval from DB (cache miss)", func(t *testing.T) {
		cacheRepo.On("Get", ctx, "leaderboard:level1:version").Return("", nil).Once()
		cacheRepo.On("Get", ctx, "leaderboard:level1:v0:1:10").Return("", nil).Once()
		repo.On("GetLegacyGlobal", ctx, levelID, 10, 0).Return([]models.LeaderboardEntry{{Rank: 1, PlayerID: "U1"}}, 1, nil).Once()
		cacheRepo.On("Set", ctx, mock.Anything, mock.Anything, 30*time.Second).Return(nil).Once()

		resp, err := usecase.GetGlobalLeaderboard(ctx, levelID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Leaderboard))
		assert.Equal(t, 1, resp.Total)
		repo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
	})

	t.Run("successful retrieval from cache", func(t *testing.T) {
		cacheRepo.On("Get", ctx, "leaderboard:level1:version").Return("5", nil).Once()

		expectedResp := models.GlobalLeaderboardResponse{
			Leaderboard: []models.LeaderboardEntry{{Rank: 1, PlayerID: "U1"}},
			Total:       1,
			Page:        1,
			PerPage:     10,
		}
		data, _ := json.Marshal(expectedResp)
		cacheRepo.On("Get", ctx, "leaderboard:level1:v5:1:10").Return(string(data), nil).Once()

		resp, err := usecase.GetGlobalLeaderboard(ctx, levelID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Leaderboard))
		assert.Equal(t, 1, resp.Total)
		cacheRepo.AssertExpectations(t)
	})
}

func TestLeaderboardUsecase_GetPlayerRank(t *testing.T) {
	repo := new(repository_mock.MockLeaderboardRepository)
	playerRepo := new(repository_mock.MockPlayerRepository)
	identityRepo := new(repository_mock.MockIdentityRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewLeaderboardUsecase(repo, playerRepo, identityRepo, cacheRepo)

	ctx := context.Background()
	userID := uuid.New()
	levelID := "level1"

	t.Run("successful retrieval with surrounding players", func(t *testing.T) {
		repo.On("GetLegacyPlayerRank", ctx, userID, levelID).Return(10, 15.5, 3, nil).Once()
		// rank 10, offset = 10 - 4 = 6
		repo.On("GetLegacyGlobal", ctx, levelID, 7, 6).Return([]models.LeaderboardEntry{
			{Rank: 7}, {Rank: 8}, {Rank: 9}, {Rank: 10}, {Rank: 11}, {Rank: 12}, {Rank: 13},
		}, 100, nil).Once()

		resp, err := usecase.GetPlayerRank(ctx, userID, levelID)

		assert.NoError(t, err)
		assert.Equal(t, 10, resp.GlobalRank)
		assert.Equal(t, 15.5, resp.BestTime)
		assert.Equal(t, 7, len(resp.SurroundingPlayers))
		repo.AssertExpectations(t)
	})
}

func TestLeaderboardUsecase_UpdateEntry(t *testing.T) {
	repo := new(repository_mock.MockLeaderboardRepository)
	playerRepo := new(repository_mock.MockPlayerRepository)
	identityRepo := new(repository_mock.MockIdentityRepository)
	cacheRepo := new(repository_mock.MockCacheRepository)
	usecase := NewLeaderboardUsecase(repo, playerRepo, identityRepo, cacheRepo)

	ctx := context.Background()
	userID := uuid.New()
	levelID := "level1"

	t.Run("successful update and cache invalidation", func(t *testing.T) {
		repo.On("UpsertLegacyEntry", ctx, userID, levelID, 12.5, 3).Return(nil).Once()
		cacheRepo.On("Get", ctx, "leaderboard:level1:version").Return("1", nil).Once()
		cacheRepo.On("Set", ctx, "leaderboard:level1:version", "2", 24*time.Hour).Return(nil).Once()

		err := usecase.UpdateEntry(ctx, userID, levelID, 12.5, 3)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
	})
}
