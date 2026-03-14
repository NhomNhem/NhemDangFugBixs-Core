package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPI_Player_GetSave(t *testing.T) {
	app, mocks := SetupTestApp()

	t.Run("successful get save", func(t *testing.T) {
		playerID, _ := uuid.Parse("test-user-id")
		expectedSave := &models.PlayerSave{
			PlayerID:    playerID,
			SaveVersion: 1,
			SaveData: models.GameSaveData{
				World: models.WorldData{Seed: 123},
			},
		}

		mocks.Player.On("GetSave", mock.Anything, playerID).Return(expectedSave, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/player/save", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var loadResp models.LoadGameResponse
		json.NewDecoder(resp.Body).Decode(&loadResp)

		assert.Equal(t, 1, loadResp.SaveVersion)
		assert.Equal(t, int64(123), loadResp.World.Seed)

		mocks.Player.AssertExpectations(t)
	})

	t.Run("save not found", func(t *testing.T) {
		playerID, _ := uuid.Parse("test-user-id")
		mocks.Player.On("GetSave", mock.Anything, playerID).Return(nil, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/player/save", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}
