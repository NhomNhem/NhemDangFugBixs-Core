package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// InitDB initializes the PostgreSQL connection pool
func InitDB() error {
	// Try DATABASE_URL first (production), fallback to SUPABASE_DATABASE_URL (local)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("SUPABASE_DATABASE_URL")
	}
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL or SUPABASE_DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse DATABASE_URL: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute
	
	// Increase connect timeout for slow networks
	config.ConnConfig.ConnectTimeout = 30 * time.Second

	// Create connection pool
	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection with longer timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("Connecting to Supabase...", "timeout", "30s")
	if err := Pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("Connected to Supabase PostgreSQL")
	return nil
}

// GetDB returns the database pool instance
func GetDB() *pgxpool.Pool {
	return Pool
}

// Close closes the database connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
		slog.Info("Database connection closed")
	}
}
