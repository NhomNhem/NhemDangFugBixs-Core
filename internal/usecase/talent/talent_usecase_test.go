package talent

import (
	"context"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	repository_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTalentUsecase_UpgradeTalent(t *testing.T) {
	talentRepo := new(repository_mock.MockTalentRepository)
	usecase := NewTalentUsecase(talentRepo)

	ctx := context.Background()
	userID := uuid.New()
	talentID := models.TalentHealth

	configs := map[string]*models.TalentConfig{
		models.TalentHealth: {
			ID:            models.TalentHealth,
			MaxLevel:      5,
			BaseCost:      100,
			CostScaling:   1.5,
			BonusPerLevel: 10,
			StatType:      "hp",
		},
	}

	t.Run("successful upgrade from level 0 to 1", func(t *testing.T) {
		talentRepo.On("GetConfigs", ctx).Return(configs, nil).Once() // For the test check if needed, but actually the usecase calls it once inside UpgradeTalent. Let's re-read usecase.

		userTalent := &models.UserTalent{
			ID:           uuid.New(),
			UserID:       userID,
			TalentID:     talentID,
			CurrentLevel: 0,
		}
		talentRepo.On("GetOrCreate", ctx, userID, talentID).Return(userTalent, nil).Once()

		// cost = 100 * (1.5 ^ 0) = 100
		talentRepo.On("UpdateUserGold", ctx, userID, -100).Return(900, nil).Once()
		talentRepo.On("UpdateLevel", ctx, userTalent.ID, 1).Return(nil).Once()

		resp, err := usecase.UpgradeTalent(ctx, userID, talentID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
		assert.Equal(t, 1, resp.NewLevel)
		assert.Equal(t, 100, resp.GoldSpent)
		assert.Equal(t, 900, resp.NewTotalGold)

		talentRepo.AssertExpectations(t)
	})

	t.Run("failed upgrade - already at max level", func(t *testing.T) {
		talentRepo.On("GetConfigs", ctx).Return(configs, nil).Once()

		userTalent := &models.UserTalent{
			ID:           uuid.New(),
			UserID:       userID,
			TalentID:     talentID,
			CurrentLevel: 5,
		}
		talentRepo.On("GetOrCreate", ctx, userID, talentID).Return(userTalent, nil).Once()

		resp, err := usecase.UpgradeTalent(ctx, userID, talentID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "talent already at max level", err.Error())
	})
}
