package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	usecase_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHollowWildsHandler_Login(t *testing.T) {
	app := fiber.New()
	authUsecase := new(usecase_mock.MockAuthUsecase)
	playerUsecase := new(usecase_mock.MockPlayerUsecase)
	analyticsUsecase := new(usecase_mock.MockAnalyticsUsecase)
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

func TestHollowWildsHandler_SaveLoad(t *testing.T) {
	app := fiber.New()
	authUsecase := new(usecase_mock.MockAuthUsecase)
	playerUsecase := new(usecase_mock.MockPlayerUsecase)
	analyticsUsecase := new(usecase_mock.MockAnalyticsUsecase)
	handler := NewHollowWildsHandler(authUsecase, playerUsecase, analyticsUsecase)

	playerID := uuid.New()

	// Mock middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userId", playerID.String())
		return c.Next()
	})

	app.Get("/player/save", handler.GetSave)
	app.Put("/player/save", handler.UpdateSave)

	t.Run("successful get save", func(t *testing.T) {
		expectedSave := &models.PlayerSave{
			PlayerID:    playerID,
			SaveVersion: 1,
			SaveData:    models.GameSaveData{World: models.WorldData{Seed: 123}},
		}
		playerUsecase.On("GetSave", mock.Anything, playerID).Return(expectedSave, nil).Once()

		req, _ := http.NewRequest("GET", "/player/save", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result models.LoadGameResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedSave.SaveVersion, result.SaveVersion)
		playerUsecase.AssertExpectations(t)
	})

	t.Run("save not found", func(t *testing.T) {
		playerUsecase.On("GetSave", mock.Anything, playerID).Return(nil, nil).Once()

		req, _ := http.NewRequest("GET", "/player/save", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
