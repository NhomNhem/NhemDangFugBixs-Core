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

	// Create user_talents table
	if err := createUserTalentsTable(ctx); err != nil {
		return err
	}

	// Create payments table
	if err := createPaymentsTable(ctx); err != nil {
		return err
	}

	// Create analytics_events table
	if err := createAnalyticsEventsTable(ctx); err != nil {
		return err
	}

	// Create Hollow Wilds tables
	if err := createHollowWildsTables(ctx); err != nil {
		return err
	}

	log.Println("✅ Migrations completed successfully")
	return nil
}

func createHollowWildsTables(ctx context.Context) error {
	query := `
	-- Players
	CREATE TABLE IF NOT EXISTS players (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		playfab_id VARCHAR(64) UNIQUE NOT NULL,
		display_name VARCHAR(64),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		last_seen_at TIMESTAMPTZ DEFAULT NOW()
	);

	-- Save Data
	CREATE TABLE IF NOT EXISTS player_saves (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		player_id UUID REFERENCES players(id) ON DELETE CASCADE,
		save_version INT DEFAULT 1,
		save_data JSONB NOT NULL,
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(player_id)
	);

	-- Save Backups
	CREATE TABLE IF NOT EXISTS player_save_backups (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		player_id UUID REFERENCES players(id) ON DELETE CASCADE,
		save_version INT,
		save_data JSONB NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	-- Leaderboard
	CREATE TABLE IF NOT EXISTS leaderboard_entries (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
		type VARCHAR(32) NOT NULL,
		value BIGINT NOT NULL,
		character VARCHAR(16) NOT NULL,
		world_seed BIGINT,
		combat_build VARCHAR(16),
		run_metadata JSONB DEFAULT '{}',
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(player_id, type, character)
	);

	CREATE INDEX IF NOT EXISTS idx_leaderboard_global ON leaderboard_entries(type, value DESC);
	CREATE INDEX IF NOT EXISTS idx_leaderboard_character ON leaderboard_entries(type, character, value DESC);
	CREATE INDEX IF NOT EXISTS idx_leaderboard_player ON leaderboard_entries(player_id);

	-- Analytics
	CREATE TABLE IF NOT EXISTS analytics_events (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID,
		event_type TEXT NOT NULL,
		event_properties JSONB DEFAULT '{}',
		session_id TEXT,
		platform TEXT,
		app_version TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_analytics_events_type ON analytics_events(event_type);
	CREATE INDEX IF NOT EXISTS idx_analytics_events_created ON analytics_events(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_analytics_events_user ON analytics_events(user_id);
	CREATE INDEX IF NOT EXISTS idx_analytics_events_session ON analytics_events(session_id);
	`

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	log.Println("  ✓ Created Hollow Wilds tables")
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

func createUserTalentsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS user_talents (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		
		-- Talent identity
		talent_id TEXT NOT NULL,
		
		-- Upgrade status
		current_level INTEGER NOT NULL DEFAULT 0 CHECK (current_level >= 0 AND current_level <= 10),
		
		-- Metadata
		upgraded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		
		UNIQUE(user_id, talent_id)
	);

	CREATE INDEX IF NOT EXISTS idx_user_talents_user ON user_talents(user_id);
	`

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	log.Println("  ✓ Created user_talents table")
	return nil
}

func createPaymentsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS payments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		
		-- Payment details
		payment_provider TEXT NOT NULL,
		provider_transaction_id TEXT UNIQUE NOT NULL,
		
		-- Purchase info
		product_id TEXT NOT NULL,
		amount_usd NUMERIC(10,2) NOT NULL CHECK (amount_usd > 0),
		diamonds_rewarded INTEGER NOT NULL CHECK (diamonds_rewarded > 0),
		
		-- Status
		status TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
		
		-- Timestamps
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		completed_at TIMESTAMPTZ,
		
		-- Metadata
		ip_address TEXT,
		user_agent TEXT,
		metadata JSONB
	);

	CREATE INDEX IF NOT EXISTS idx_payments_user ON payments(user_id);
	CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
	CREATE INDEX IF NOT EXISTS idx_payments_created ON payments(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_payments_provider_tx ON payments(provider_transaction_id);
	`

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	log.Println("  ✓ Created payments table")
	return nil
}

func createAnalyticsEventsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS analytics_events (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		
		-- Event details
		event_type TEXT NOT NULL,
		event_properties JSONB,
		
		-- Context
		session_id TEXT,
		platform TEXT,
		app_version TEXT,
		
		-- Timestamp
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_analytics_events_user ON analytics_events(user_id);
	CREATE INDEX IF NOT EXISTS idx_analytics_events_type ON analytics_events(event_type);
	CREATE INDEX IF NOT EXISTS idx_analytics_events_created ON analytics_events(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_analytics_events_session ON analytics_events(session_id);
	`

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	log.Println("  ✓ Created analytics_events table")
	return nil
}
