package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "github.com/NhomNhem/HollowWilds-Backend/docs"
	"github.com/NhomNhem/HollowWilds-Backend/internal/database"
	"github.com/NhomNhem/HollowWilds-Backend/internal/delivery/http"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/infrastructure/cache"
	"github.com/NhomNhem/HollowWilds-Backend/internal/infrastructure/identity"
	"github.com/NhomNhem/HollowWilds-Backend/internal/infrastructure/persistence"
	"github.com/NhomNhem/HollowWilds-Backend/internal/middleware"
	"github.com/NhomNhem/HollowWilds-Backend/internal/usecase/admin"
	"github.com/NhomNhem/HollowWilds-Backend/internal/usecase/analytics"
	"github.com/NhomNhem/HollowWilds-Backend/internal/usecase/auth"
	"github.com/NhomNhem/HollowWilds-Backend/internal/usecase/leaderboard"
	// "github.com/NhomNhem/HollowWilds-Backend/internal/usecase/level" // DISABLED
	"github.com/NhomNhem/HollowWilds-Backend/internal/usecase/player"
	// "github.com/NhomNhem/HollowWilds-Backend/internal/usecase/talent" // DISABLED
	"github.com/NhomNhem/HollowWilds-Backend/pkg/utils"
)

// @title Hollow Wilds Backend API
// @version 1.1.0
// @description Game backend API với PlayFab integration, anti-cheat validation, và talent system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@gamefeel.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host gamefeel-backend.fly.dev
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token format: "Bearer {token}"

// @securityDefinitions.apikey PlayFabToken
// @in header
// @name X-PlayFab-SessionToken
// @description PlayFab session token for authentication
var (
	Version = "dev"
	Commit  = "none"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize structured logging
	var handler slog.Handler
	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(handler))

	slog.Info("Starting server initialization", "version", Version, "commit", Commit, "env", getEnv("ENV", "development"))

	// Security Hardening: Validate environment
	isProd := os.Getenv("ENV") == "production"
	jwtSecret := os.Getenv("JWT_SECRET")
	if isProd {
		if jwtSecret == "" || len(jwtSecret) < 32 {
			slog.Error("CRITICAL: JWT_SECRET is too short or missing in production! System must have at least 32 characters.")
			os.Exit(1)
		}
		if os.Getenv("PLAYFAB_TITLE_ID") == "" || os.Getenv("PLAYFAB_TITLE_ID") == "DEV" {
			slog.Error("CRITICAL: PLAYFAB_TITLE_ID cannot be empty or 'DEV' in production!")
			os.Exit(1)
		}
		if os.Getenv("ALLOWED_ORIGINS") == "" || os.Getenv("ALLOWED_ORIGINS") == "*" {
			slog.Warn("SECURITY: ALLOWED_ORIGINS is not set or permissive in production. Consider restricting this.")
		}
	}

	// Initialize database connection (optional for development)
	if err := database.InitDB(); err != nil {
		slog.Error("Database connection failed", "error", err)
		slog.Warn("Continuing without database (API endpoints will return mock data)")
	}
	defer database.Close()

	// Initialize Redis connection
	if err := utils.InitRedis(); err != nil {
		slog.Error("Redis connection failed", "error", err)
		slog.Warn("Continuing without Redis (caching and sessions will be disabled)")
	}
	defer utils.CloseRedis()

	// Create Redis storage for Fiber middleware
	var fiberRedisStorage fiber.Storage
	if utils.RedisClient != nil {
		fiberRedisStorage = redis.New(redis.Config{
			URL: os.Getenv("UPSTASH_REDIS_URL"),
		})
	}

	// Run automated database migrations (only if connected)
	if database.Pool != nil {
		if err := database.RunMigrations(); err != nil {
			slog.Error("Automated migrations failed", "error", err)
			slog.Warn("Continuing server startup despite migration failure")
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Hollow Wilds Backend v1.1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			requestID, _ := c.Locals("requestId").(string)

			return c.Status(code).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    getSnakeCaseErrorCode(code),
					Message: err.Error(),
					TraceID: requestID,
				},
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggerMiddleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-PlayFab-SessionToken",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Global rate limiter (basic protection)
	app.Use(limiter.New(limiter.Config{
		Storage:    fiberRedisStorage,
		Max:        getEnvInt("RATE_LIMIT_REQUESTS", 100),
		Expiration: getEnvDuration("RATE_LIMIT_DURATION", 60*time.Second),
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
				},
			})
		},
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":  "ok",
			"message": "Hollow Wilds Backend is running",
			"version": "1.1.0",
		}

		// Check database connection (optional)
		if database.Pool != nil {
			ctx := c.Context()
			if err := database.Pool.Ping(ctx); err != nil {
				health["database"] = "disconnected"
				health["database_error"] = err.Error()
			} else {
				health["database"] = "connected"
			}
		} else {
			health["database"] = "not configured"
		}

		return c.JSON(health)
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1 routes
	apiV1 := app.Group("/api/v1")

	// Root endpoint with API info
	apiV1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "GameFeel API v1 (Hollow Wilds Phase 1)",
			"endpoints": []string{
				"GET  /health",
				"GET  /api/v1/",
				"POST /api/v1/auth/hw/login",
				"PUT  /api/v1/player/save",
				"GET  /api/v1/leaderboard",
				"POST /api/v1/analytics/events",
			},
		})
	})

	// Initialize Infrastructure
	playerRepo := persistence.NewPostgresPlayerRepository(database.Pool)
	saveRepo := persistence.NewPostgresSaveRepository(database.Pool)
	leaderboardRepo := persistence.NewPostgresLeaderboardRepository(database.Pool)
	analyticsRepo := persistence.NewPostgresAnalyticsRepository(database.Pool)
	// talentRepo := persistence.NewPostgresTalentRepository(database.Pool) // DISABLED
	// levelRepo := persistence.NewPostgresLevelRepository(database.Pool) // DISABLED
	adminRepo := persistence.NewPostgresAdminRepository(database.Pool)
	redisRepo := cache.NewRedisRepository(utils.RedisClient)
	identityRepo := identity.NewPlayFabRepository()

	// Initialize Usecases
	authUsecase := auth.NewAuthUsecase(playerRepo, identityRepo, redisRepo)
	playerUsecase := player.NewPlayerUsecase(playerRepo, saveRepo, redisRepo)
	leaderboardUsecase := leaderboard.NewLeaderboardUsecase(leaderboardRepo, playerRepo, identityRepo, redisRepo)
	analyticsUsecase := analytics.NewAnalyticsUsecase(analyticsRepo)
	// talentUsecase := talent.NewTalentUsecase(talentRepo) // DISABLED
	// levelUsecase := level.NewLevelUsecase(levelRepo, leaderboardUsecase) // DISABLED
	adminUsecase := admin.NewAdminUsecase(adminRepo, leaderboardRepo)

	// Register handlers
	// authHandler := http.NewAuthHandler(authUsecase) // DISABLED
	// levelHandler := http.NewLevelHandler(levelUsecase) // DISABLED
	// talentHandler := http.NewTalentHandler(talentUsecase) // DISABLED
	leaderboardHandler := http.NewLeaderboardHandler(leaderboardUsecase)
	adminHandler := http.NewAdminHandler(adminUsecase)
	hollowWildsHandler := http.NewHollowWildsHandler(authUsecase, playerUsecase, analyticsUsecase)

	// Auth routes (public)
	auth := apiV1.Group("/auth")
	// auth.Post("/login", authHandler.Login) // LEGACY DISABLED
	auth.Post("/hw/login", hollowWildsHandler.Login)
	auth.Post("/refresh", hollowWildsHandler.Refresh)
	auth.Delete("/logout", hollowWildsHandler.Logout)

	// Player routes (Hollow Wilds)
	player := apiV1.Group("/player", middleware.AuthMiddleware())
	player.Get("/save", hollowWildsHandler.GetSave)
	player.Put("/save", hollowWildsHandler.UpdateSave)
	player.Post("/save/backup", hollowWildsHandler.CreateBackup)
	player.Get("/save/backups", hollowWildsHandler.GetBackups)
	player.Post("/save/restore", hollowWildsHandler.RestoreFromBackup)

	// Protected routes (require JWT) - DISABLED LEGACY
	/*
		levels := apiV1.Group("/levels", middleware.AuthMiddleware())
		levels.Post("/complete", levelHandler.CompleteLevel)

		talents := apiV1.Group("/talents", middleware.AuthMiddleware())
		talents.Get("/", talentHandler.GetTalents)
		talents.Post("/upgrade", talentHandler.UpgradeTalent)
	*/

	// Leaderboard routes
	lbHw := apiV1.Group("/leaderboard")
	lbHw.Get("/", leaderboardHandler.GetHollowWildsLeaderboard)
	lbHw.Post("/submit", middleware.AuthMiddleware(), leaderboardHandler.SubmitHollowWildsEntry)
	lbHw.Get("/player", middleware.AuthMiddleware(), leaderboardHandler.GetPlayerHollowWildsStats)

	// Legacy Leaderboard routes - DISABLED
	/*
		leaderboards := apiV1.Group("/leaderboards")
		leaderboards.Get("/:levelId", leaderboardHandler.GetGlobalLeaderboard)
		leaderboards.Get("/:levelId/me", middleware.AuthMiddleware(), leaderboardHandler.GetPlayerRank)
		leaderboards.Get("/:levelId/friends", middleware.AuthMiddleware(), leaderboardHandler.GetFriendsLeaderboard)
	*/

	// Analytics routes
	analytics := apiV1.Group("/analytics")
	analytics.Post("/events", limiter.New(limiter.Config{
		Storage:    fiberRedisStorage,
		Max:        100,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}), middleware.AuthMiddleware(), hollowWildsHandler.TrackEvents)
	// Admin routes (require JWT + admin role)
	admin := apiV1.Group("/admin", middleware.AuthMiddleware(), middleware.AdminMiddleware())
	admin.Get("/users/search", adminHandler.SearchUsers)
	admin.Get("/users/:userId/profile", adminHandler.GetUserProfile)
	admin.Post("/users/:userId/adjust-gold", adminHandler.AdjustGold)
	admin.Post("/users/:userId/ban", adminHandler.BanUser)
	admin.Post("/users/:userId/unban", adminHandler.UnbanUser)
	admin.Get("/users/:userId/ban-history", adminHandler.GetBanHistory)
	admin.Get("/users/:userId/export-data", adminHandler.ExportUserData)
	admin.Get("/actions", adminHandler.GetAdminActions)
	admin.Get("/stats/overview", adminHandler.GetSystemStats)
	admin.Delete("/leaderboards/:levelId", adminHandler.ResetLeaderboard)
	admin.Get("/leaderboards/stats", adminHandler.GetLeaderboardStats)

	// Get port from env or default to 8080
	port := getEnv("PORT", "8080")

	// Start server in a goroutine
	go func() {
		slog.Info("Server starting", "port", port)
		if err := app.Listen(":" + port); err != nil {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutdown signal received, initiating graceful shutdown...")

	// Create shutdown context with 10s timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Stop Fiber (stop accepting new requests)
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		slog.Error("Fiber shutdown failed", "error", err)
	} else {
		slog.Info("Fiber stopped successfully")
	}

	// 2. Close Database pool
	database.Close()

	// 3. Close Redis connection
	utils.CloseRedis()

	slog.Info("Graceful shutdown complete. Exiting.")
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as int with default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// Helper function to get environment variable as duration with default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getSnakeCaseErrorCode(code int) string {
	switch code {
	case fiber.StatusBadRequest:
		return "invalid_request"
	case fiber.StatusUnauthorized:
		return "unauthorized"
	case fiber.StatusForbidden:
		return "forbidden"
	case fiber.StatusNotFound:
		return "not_found"
	case fiber.StatusMethodNotAllowed:
		return "method_not_allowed"
	case fiber.StatusConflict:
		return "conflict"
	case fiber.StatusUnprocessableEntity:
		return "validation_error"
	case fiber.StatusTooManyRequests:
		return "rate_limited"
	case fiber.StatusInternalServerError:
		return "internal_error"
	default:
		return "error"
	}
}
