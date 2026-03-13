# Supabase Setup Guide

**Task**: Phase 1.3 - Setup PostgreSQL Database  
**Provider**: Supabase (PostgreSQL 15)  
**Last Updated**: 2026-03-04

---

## Step 1: Create Supabase Account

### 1.1 Sign Up

Go to: https://supabase.com

**Sign up with:**
- GitHub account (recommended - fastest)
- Or email/password

### 1.2 Create New Project

1. Click **"New Project"**
2. Fill in details:
   - **Name**: `gamefeel-backend` (or any name)
   - **Database Password**: Generate strong password (SAVE THIS!)
   - **Region**: Choose closest to your users (e.g., `Southeast Asia (Singapore)` or `US East`)
   - **Pricing Plan**: **Free** (500MB database, 2GB bandwidth)

3. Click **"Create new project"**
4. Wait ~2 minutes for provisioning

---

## Step 2: Get Connection Details

### 2.1 Database Connection String

1. In Supabase dashboard, go to **Settings** (gear icon)
2. Click **Database** in left sidebar
3. Scroll to **Connection string** section
4. Select **"URI"** tab
5. Copy the connection string:

```
postgresql://postgres.[project-ref]:[YOUR-PASSWORD]@aws-0-[region].pooler.supabase.com:6543/postgres
```

**IMPORTANT**: Replace `[YOUR-PASSWORD]` with the password you saved!

### 2.2 Service Role Key

1. Go to **Settings** → **API**
2. Scroll to **Project API keys**
3. Copy **`service_role`** key (secret key - NEVER expose in client!)

**Example:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFiY2RlZiIsInJvbGUiOiJzZXJ2aWNlX3JvbGUiLCJpYXQiOjE2Nzg5ODcyMDAsImV4cCI6MTk5NDU2MzIwMH0.abcdef123456...
```

### 2.3 Anon Key (for Unity client - later)

Also copy **`anon`** key (public key - safe to expose):
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFiY2RlZiIsInJvbGUiOiJhbm9uIiwiaWF0IjoxNjc4OTg3MjAwLCJleHAiOjE5OTQ1NjMyMDB9.abcdef123456...
```

---

## Step 3: Update .env File

Open `I:\unityVers\GameFeel-Backend\.env` and update:

```env
# Server Configuration
PORT=8080
ENV=development

# Database (Supabase) - UPDATE THESE!
SUPABASE_DATABASE_URL=postgresql://postgres.[project-ref]:[YOUR-PASSWORD]@aws-0-[region].pooler.supabase.com:6543/postgres
SUPABASE_SERVICE_ROLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Supabase Project Details (for reference)
SUPABASE_PROJECT_URL=https://[project-ref].supabase.co
SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# JWT Secret (Change in production!)
JWT_SECRET=your-secret-key-here-change-in-production

# CORS
ALLOWED_ORIGINS=*

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=60s
```

**Save the file!**

---

## Step 4: Install PostgreSQL Driver

```powershell
cd I:\unityVers\GameFeel-Backend

# Install pgx (PostgreSQL driver)
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
```

---

## Step 5: Create Database Connection Code

Create `internal/database/db.go`:

```go
package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// InitDB initializes the PostgreSQL connection pool
func InitDB() error {
	dbURL := os.Getenv("SUPABASE_DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("SUPABASE_DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse DATABASE_URL: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute

	// Create connection pool
	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✅ Connected to Supabase PostgreSQL")
	return nil
}

// Close closes the database connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("🔌 Database connection closed")
	}
}
```

---

## Step 6: Update main.go to Connect Database

Update `cmd/server/main.go`:

```go
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/yourusername/GameFeel-Backend/internal/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Hollow Wilds Backend v1.0.0",
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
		// Check database connection
		ctx := c.Context()
		if err := database.Pool.Ping(ctx); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":  "error",
				"message": "Database connection failed",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"message":  "Hollow Wilds Backend is running",
			"version":  "1.0.0",
			"database": "connected",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	api.Get("/", func(c *fiber.Ctx) error {
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

---

## Step 7: Test Connection

```powershell
cd I:\unityVers\GameFeel-Backend

# Run server
go run cmd/server/main.go

# Expected output:
# ✅ Connected to Supabase PostgreSQL
# 🚀 Server starting on port 8080...
```

**Test health check:**
```powershell
curl http://localhost:8080/health

# Should return:
# {
#   "status": "ok",
#   "message": "Hollow Wilds Backend is running",
#   "version": "1.0.0",
#   "database": "connected"
# }
```

---

## Step 8: Create Database Tables (Migrations)

We'll create tables from `database-schema.md`.

Create `internal/database/migrations.go`:

```go
package database

import (
	"context"
	"log"
)

// RunMigrations creates all necessary database tables
func RunMigrations() error {
	ctx := context.Background()

	log.Println("🔄 Running database migrations...")

	// Create users table
	if err := createUsersTable(ctx); err != nil {
		return err
	}

	// Create level_completions table
	if err := createLevelCompletionsTable(ctx); err != nil {
		return err
	}

	// Add more tables as needed...

	log.Println("✅ Migrations completed successfully")
	return nil
}

func createUsersTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		playfab_id TEXT UNIQUE NOT NULL,
		display_name TEXT,
		
		-- Currency
		gold INTEGER NOT NULL DEFAULT 0 CHECK (gold >= 0),
		diamonds INTEGER NOT NULL DEFAULT 0 CHECK (diamonds >= 0),
		
		-- Progression
		max_map_unlocked INTEGER NOT NULL DEFAULT 1 CHECK (max_map_unlocked >= 1),
		total_stars_collected INTEGER NOT NULL DEFAULT 0 CHECK (total_stars_collected >= 0),
		
		-- Metadata
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		last_login_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		last_played_at TIMESTAMPTZ,
		total_play_time_seconds INTEGER DEFAULT 0,
		
		-- Social
		facebook_id TEXT UNIQUE,
		google_id TEXT UNIQUE,
		
		-- Flags
		is_banned BOOLEAN DEFAULT FALSE,
		ban_reason TEXT,
		banned_at TIMESTAMPTZ,
		
		-- GDPR
		deleted_at TIMESTAMPTZ
	);

	CREATE INDEX IF NOT EXISTS idx_users_playfab_id ON users(playfab_id);
	CREATE INDEX IF NOT EXISTS idx_users_last_login ON users(last_login_at DESC);
	CREATE INDEX IF NOT EXISTS idx_users_active ON users(last_played_at DESC) WHERE deleted_at IS NULL;
	`

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	log.Println("  ✓ Created users table")
	return nil
}

func createLevelCompletionsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS level_completions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		
		-- Level identity
		level_id TEXT NOT NULL,
		map_id TEXT NOT NULL,
		
		-- Performance stats
		stars_earned INTEGER NOT NULL CHECK (stars_earned BETWEEN 0 AND 3),
		best_time_seconds REAL NOT NULL CHECK (best_time_seconds > 0),
		play_count INTEGER NOT NULL DEFAULT 1,
		
		-- Latest run stats
		last_final_hp REAL,
		last_dash_count INTEGER,
		last_counter_count INTEGER,
		last_vulnerable_kills INTEGER,
		
		-- Timestamps
		first_completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		last_played_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		
		UNIQUE(user_id, level_id)
	);

	CREATE INDEX IF NOT EXISTS idx_level_completions_user ON level_completions(user_id);
	CREATE INDEX IF NOT EXISTS idx_level_completions_level ON level_completions(level_id);
	CREATE INDEX IF NOT EXISTS idx_level_completions_stars ON level_completions(user_id, stars_earned);
	`

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	log.Println("  ✓ Created level_completions table")
	return nil
}
```

**Update main.go to run migrations:**

```go
// After database.InitDB()
if err := database.RunMigrations(); err != nil {
	log.Printf("⚠️  Migration failed: %v", err)
}
```

---

## Step 9: Verify Tables in Supabase Dashboard

1. Go to Supabase dashboard
2. Click **Table Editor** in left sidebar
3. You should see:
   - `users` table
   - `level_completions` table

---

## Troubleshooting

### Error: "unable to ping database"

**Check:**
- Is `.env` file updated with correct credentials?
- Is password correct? (no special characters escaped?)
- Is Supabase project running? (check dashboard)

### Error: "connection refused"

**Try:**
- Use **Session mode** connection string (port 5432) instead of Pooler (port 6543)
- Check firewall settings

### Error: "permission denied"

**Solution:**
- Make sure you're using **service_role** key, not anon key
- Check that RLS is disabled for migrations (we'll enable it later)

---

## Next Steps

After Supabase is connected:
- ✅ Database connection working
- ✅ Basic tables created
- → Setup Redis cache (Task 1.4)
- → Implement auth endpoint (Task 1.5)

---

**Status**: Ready to connect Supabase!  
**Let me know when you have your Supabase credentials, and I'll help you set it up!**
