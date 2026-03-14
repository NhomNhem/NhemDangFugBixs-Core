package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPI_Auth_Login(t *testing.T) {
	app, mocks := SetupTestApp()

	t.Run("successful login", func(t *testing.T) {
		reqBody := models.AuthRequest{
			PlayFabID: "PLAYFAB_123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mocks.Auth.On("LegacyLogin", mock.Anything, "PLAYFAB_123", "", "valid-token").Return(&models.AuthResponse{
			JWT: "jwt-token",
			User: models.User{
				PlayFabID: "PLAYFAB_123",
			},
			ExpiresIn: 3600,
		}, nil).Once()

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-PlayFab-SessionToken", "valid-token")

		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var apiResp models.APIResponse
		json.NewDecoder(resp.Body).Decode(&apiResp)

		assert.True(t, apiResp.Success)
		data := apiResp.Data.(map[string]interface{})
		assert.Equal(t, "jwt-token", data["jwt"])

		mocks.Auth.AssertExpectations(t)
	})
}
