package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/mocks/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHollowWildsHandler_Login(t *testing.T) {
	app := fiber.New()
	authUsecase := new(usecase_mock.AuthUsecase)
	playerUsecase := new(usecase_mock.PlayerUsecase)
	analyticsUsecase := new(usecase_mock.AnalyticsUsecase)
	handler := NewHollowWildsHandler(authUsecase, playerUsecase, analyticsUsecase)

	app.Post("/auth/hw/login", handler.Login)

	t.Run("successful login", func(t *testing.T) {
		reqBody := models.HollowWildsLoginRequest{
			PlayfabSessionTicket: "valid-ticket",
		}
		expectedResp := &models.HollowWildsAuthResponse{
			Token:    "jwt-token",
			PlayerID: "player-123",
		}

		authUsecase.On("Login", mock.Anything, reqBody.PlayfabSessionTicket, "PF_123").Return(expectedResp, nil).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/hw/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-PlayFab-ID", "PF_123")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result models.HollowWildsAuthResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResp.Token, result.Token)
		authUsecase.AssertExpectations(t)
	})
}
