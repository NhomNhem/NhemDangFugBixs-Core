package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPI_Leaderboard_Get(t *testing.T) {
	app, mocks := SetupTestApp()

	t.Run("successful get leaderboard", func(t *testing.T) {
		expectedResp := &models.HollowWildsLeaderboardResponse{
			Type: "longest_run_days",
			Entries: []models.HollowWildsLeaderboardEntry{
				{Rank: 1, PlayerID: "p1", Value: 100},
			},
		}

		mocks.Leaderboard.On("GetLeaderboard", mock.Anything, "longest_run_days", "global", "", 100, 0).Return(expectedResp, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/leaderboard?type=longest_run_days", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var lbResp models.HollowWildsLeaderboardResponse
		json.NewDecoder(resp.Body).Decode(&lbResp)

		assert.Equal(t, "longest_run_days", lbResp.Type)
		assert.Equal(t, 1, len(lbResp.Entries))

		mocks.Leaderboard.AssertExpectations(t)
	})
}
