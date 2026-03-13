package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/delivery/http"
	usecase_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// SetupTestApp creates a Fiber app instance for testing with mocked usecases
func SetupTestApp() (*fiber.App, *usecase_mock.MockAuthUsecase) {
	app := fiber.New()

	mockAuthUsecase := new(usecase_mock.MockAuthUsecase)
	authHandler := http.NewAuthHandler(mockAuthUsecase)

	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
	})
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	return app, mockAuthUsecase
}

func TestAPI_HealthCheck(t *testing.T) {
	app, _ := SetupTestApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "ok", body["status"])
}

func TestAPI_RootV1(t *testing.T) {
	app, _ := SetupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/", nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
