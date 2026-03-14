package api

import (
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/delivery/http"
	usecase_mock "github.com/NhomNhem/HollowWilds-Backend/internal/mocks/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// TestMocks holds all mocked usecases
type TestMocks struct {
	Auth        *usecase_mock.MockAuthUsecase
	Player      *usecase_mock.MockPlayerUsecase
	Leaderboard *usecase_mock.MockLeaderboardUsecase
	Analytics   *usecase_mock.MockAnalyticsUsecase
}

// SetupTestApp creates a Fiber app instance for testing with mocked usecases
func SetupTestApp() (*fiber.App, *TestMocks) {
	app := fiber.New()

	// Add basic limiter for testing
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
	}))

	mocks := &TestMocks{
		Auth:        new(usecase_mock.MockAuthUsecase),
		Player:      new(usecase_mock.MockPlayerUsecase),
		Leaderboard: new(usecase_mock.MockLeaderboardUsecase),
		Analytics:   new(usecase_mock.MockAnalyticsUsecase),
	}

	authHandler := http.NewAuthHandler(mocks.Auth)
	hollowWildsHandler := http.NewHollowWildsHandler(mocks.Auth, mocks.Player, mocks.Analytics)
	leaderboardHandler := http.NewLeaderboardHandler(mocks.Leaderboard)

	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
	})

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/hw/login", hollowWildsHandler.Login)

	// Player (mocking auth middleware by using a group without it for simple integration tests,
	// or we can add a dummy middleware that sets locals)
	player := api.Group("/player")
	player.Use(func(c *fiber.Ctx) error {
		c.Locals("userId", "test-user-id")
		return c.Next()
	})
	player.Get("/save", hollowWildsHandler.GetSave)
	player.Put("/save", hollowWildsHandler.UpdateSave)

	// Leaderboard
	api.Get("/leaderboard", leaderboardHandler.GetHollowWildsLeaderboard)

	// Analytics
	api.Post("/analytics/events", hollowWildsHandler.TrackEvents)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	return app, mocks
}
