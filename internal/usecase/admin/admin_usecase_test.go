package admin

import (
	"context"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	repository_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminUsecase_AdjustGold(t *testing.T) {
	adminRepo := new(repository_mock.MockAdminRepository)
	lbRepo := new(repository_mock.MockLeaderboardRepository)
	usecase := NewAdminUsecase(adminRepo, lbRepo)

	ctx := context.Background()
	adminID := uuid.New()
	userID := uuid.New()
	amount := 500
	reason := "Customer support reward"
	ip := "127.0.0.1"

	t.Run("successful gold adjustment", func(t *testing.T) {
		expectedResp := &models.AdjustGoldResponse{
			UserID:     userID,
			OldBalance: 1000,
			NewBalance: 1500,
			Adjustment: amount,
			Reason:     reason,
		}
		adminRepo.On("AdjustGold", ctx, adminID, userID, amount, reason, ip).Return(expectedResp, nil).Once()

		resp, err := usecase.AdjustGold(ctx, adminID, userID, amount, reason, ip)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
		adminRepo.AssertExpectations(t)
	})
}

func TestAdminUsecase_ResetLeaderboard(t *testing.T) {
	adminRepo := new(repository_mock.MockAdminRepository)
	lbRepo := new(repository_mock.MockLeaderboardRepository)
	usecase := NewAdminUsecase(adminRepo, lbRepo)

	ctx := context.Background()
	adminID := uuid.New()
	levelID := "level1"
	reason := "Cheating cleanup"
	ip := "127.0.0.1"

	t.Run("successful reset", func(t *testing.T) {
		lbRepo.On("ResetLegacyLeaderboard", ctx, levelID).Return(nil).Once()
		adminRepo.On("LogAction", ctx, adminID, "RESET_LEADERBOARD", mock.Anything, mock.Anything, ip).Return(nil).Once()

		err := usecase.ResetLeaderboard(ctx, adminID, levelID, reason, ip)

		assert.NoError(t, err)
		lbRepo.AssertExpectations(t)
		adminRepo.AssertExpectations(t)
	})
}

func TestAdminUsecase_GetLeaderboardStats(t *testing.T) {
	adminRepo := new(repository_mock.MockAdminRepository)
	lbRepo := new(repository_mock.MockLeaderboardRepository)
	usecase := NewAdminUsecase(adminRepo, lbRepo)

	ctx := context.Background()

	t.Run("successful stats retrieval", func(t *testing.T) {
		expectedResp := &models.LeaderboardStatsResponse{
			TotalEntries: 100,
		}
		lbRepo.On("GetLegacyStats", ctx).Return(expectedResp, nil).Once()

		resp, err := usecase.GetLeaderboardStats(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
		lbRepo.AssertExpectations(t)
	})
}
