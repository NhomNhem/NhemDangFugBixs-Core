package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
)

// RunMigrations scans the migrations directory and applies any pending SQL files
func RunMigrations() error {
	ctx := context.Background()

	slog.Info("Starting automated database migrations...")

	// 1. Acquire advisory lock to prevent concurrent migrations
	// Using a random 64-bit integer as the lock ID
	const lockID = 123456789
	_, err := Pool.Exec(ctx, "SELECT pg_advisory_lock($1)", lockID)
	if err != nil {
		return fmt.Errorf("failed to acquire migration lock: %w", err)
	}
	defer Pool.Exec(ctx, "SELECT pg_advisory_unlock($1)", lockID)

	// 2. Create tracking table if it doesn't exist
	_, err = Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// 3. Scan migrations directory
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// 4. Apply pending migrations
	appliedCount := 0
	for _, file := range files {
		var exists bool
		err := Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", file).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", file, err)
		}

		if exists {
			continue
		}

		slog.Info("Applying migration", "file", file)
		content, err := os.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute migration in a transaction
		tx, err := Pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", file, err)
		}
		defer tx.Rollback(ctx)

		if _, err := tx.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", file); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		appliedCount++
		slog.Info("Successfully applied migration", "file", file)
	}

	if appliedCount == 0 {
		slog.Info("Database is up to date. No pending migrations.")
	} else {
		slog.Info("Finished applying migrations", "count", appliedCount)
	}

	return nil
}
