-- Migration: Legacy Leaderboards
-- Description: Adds table for level-based rankings
-- Date: 2026-03-14

CREATE TABLE IF NOT EXISTS level_rankings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id TEXT NOT NULL,
    best_time_seconds REAL NOT NULL,
    stars_earned INTEGER NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, level_id)
);

-- Index for global ranking per level
CREATE INDEX IF NOT EXISTS idx_level_rankings_score ON level_rankings(level_id, best_time_seconds ASC);

-- Index for player lookups
CREATE INDEX IF NOT EXISTS idx_level_rankings_user ON level_rankings(user_id);

-- 1.3 Create leaderboard_cache table for Redis sync tracking (optional)
CREATE TABLE IF NOT EXISTS leaderboard_cache (
    level_id TEXT PRIMARY KEY,
    last_synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    needs_refresh BOOLEAN DEFAULT FALSE
);
