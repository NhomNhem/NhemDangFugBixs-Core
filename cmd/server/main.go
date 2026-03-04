package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/NhomNhem/GameFeel-Backend/internal/api"
	"github.com/NhomNhem/GameFeel-Backend/internal/database"
	"github.com/NhomNhem/GameFeel-Backend/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection (optional for development)
	if err := database.InitDB(); err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Println("⚠️  Continuing without database (API endpoints will return mock data)")
		log.Println("💡 To fix: Check network/firewall, or deploy to server with IPv6 support")
	}
	defer database.Close()

	// Run database migrations (only if connected)
	if database.Pool != nil {
		if err := database.RunMigrations(); err != nil {
			log.Printf("⚠️  Migration failed: %v", err)
			log.Println("Continuing without migrations (tables may already exist)")
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "GameFeel Backend v1.0.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    code,
					"message": err.Error(),
				},
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":  "ok",
			"message": "GameFeel Backend is running",
			"version": "1.0.0",
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

	// API v1 routes
	apiV1 := app.Group("/api/v1")

	// Root endpoint with API info
	apiV1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "GameFeel API v1",
			"endpoints": []string{
				"GET  /health",
				"GET  /api/v1/",
				"POST /api/v1/auth/login",
				"POST /api/v1/levels/complete",
				"POST /api/v1/talents/upgrade",
				"POST /api/v1/payments/create-session",
				"POST /api/v1/analytics/events",
			},
		})
	})

	// Register handlers
	authHandler := api.NewAuthHandler()
	levelHandler := api.NewLevelHandler()
	talentHandler := api.NewTalentHandler()
	
	// Auth routes (public)
	auth := apiV1.Group("/auth")
	auth.Post("/login", authHandler.Login)
	
	// Protected routes (require JWT)
	levels := apiV1.Group("/levels", middleware.AuthMiddleware())
	levels.Post("/complete", levelHandler.CompleteLevel)
	
	talents := apiV1.Group("/talents", middleware.AuthMiddleware())
	talents.Get("/", talentHandler.GetTalents)
	talents.Post("/upgrade", talentHandler.UpgradeTalent)

	// Get port from env or default to 8080
	port := getEnv("PORT", "8080")

	// Start server
	log.Printf("🚀 Server starting on port %s...", port)
	log.Printf("📝 Environment: %s", getEnv("ENV", "development"))
	log.Printf("🔗 Health check: http://localhost:%s/health", port)
	log.Printf("🔗 API docs: http://localhost:%s/api/v1/", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
