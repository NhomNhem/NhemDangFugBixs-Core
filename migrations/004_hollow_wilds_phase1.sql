-- Migration: Hollow Wilds Phase 1
-- Description: Adds players, player_saves, and player_save_backups tables for Hollow Wilds game
-- Date: 2026-03-12

-- ============================================================================
-- 1. Create players table
-- ============================================================================
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    playfab_id VARCHAR(64) UNIQUE NOT NULL,
    display_name VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index on playfab_id for fast lookups (unique constraint already creates index)
CREATE INDEX IF NOT EXISTS idx_players_playfab_id ON players(playfab_id);

-- Index on last_seen_at for activity queries
CREATE INDEX IF NOT EXISTS idx_players_last_seen_at ON players(last_seen_at DESC);

-- ============================================================================
-- 2. Create player_saves table
-- ============================================================================
CREATE TABLE IF NOT EXISTS player_saves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    save_version INT DEFAULT 1 NOT NULL,
    save_data JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(player_id)
);

-- Unique index on player_id for fast lookups
CREATE INDEX IF NOT EXISTS idx_player_saves_player_id ON player_saves(player_id);

-- Index on updated_at for tracking recent saves
CREATE INDEX IF NOT EXISTS idx_player_saves_updated_at ON player_saves(updated_at DESC);

-- ============================================================================
-- 3. Create player_save_backups table
-- ============================================================================
CREATE TABLE IF NOT EXISTS player_save_backups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    save_version INT NOT NULL,
    save_data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index on player_id for listing backups
CREATE INDEX IF NOT EXISTS idx_player_save_backups_player_id ON player_save_backups(player_id);

-- Index on created_at for ordering backups
CREATE INDEX IF NOT EXISTS idx_player_save_backups_created_at ON player_save_backups(created_at DESC);

-- Composite index for player-specific backup queries
CREATE INDEX IF NOT EXISTS idx_player_save_backups_player_created ON player_save_backups(player_id, created_at DESC);

-- ============================================================================
-- 4. Analytics events table (Shared between legacy and Hollow Wilds)
-- ============================================================================
CREATE TABLE IF NOT EXISTS analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,                      -- References users(id) OR players(id)
    event_type TEXT NOT NULL,          -- Renamed from event_name to match context
    event_properties JSONB DEFAULT '{}', -- Renamed from payload to match context
    session_id TEXT,
    platform TEXT,
    app_version TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index on event_type for analytics queries
CREATE INDEX IF NOT EXISTS idx_analytics_events_type ON analytics_events(event_type);

-- Index on created_at for time-based queries
CREATE INDEX IF NOT EXISTS idx_analytics_events_created ON analytics_events(created_at DESC);

-- Index on user_id for player-specific analytics
CREATE INDEX IF NOT EXISTS idx_analytics_events_user ON analytics_events(user_id);

-- Index on session_id for session tracking
CREATE INDEX IF NOT EXISTS idx_analytics_events_session ON analytics_events(session_id);

-- ============================================================================
-- 5. Create leaderboard_entries table
-- ============================================================================
-- type: longest_run_days | sebilah_soul_level | bosses_killed
-- Stores personal best per (player, type, character)
CREATE TABLE IF NOT EXISTS leaderboard_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  type VARCHAR(32) NOT NULL,
  value BIGINT NOT NULL,
  character VARCHAR(16) NOT NULL,    -- RIMBA | DARA | BAYU | SARI
  world_seed BIGINT,
  combat_build VARCHAR(16),          -- balanced | berserker | shade_walker
  run_metadata JSONB DEFAULT '{}',   -- bosses_killed, biomes_explored, etc
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(player_id, type, character) -- 1 personal best per player+type+character
);

-- Global leaderboard query
CREATE INDEX IF NOT EXISTS idx_leaderboard_global ON leaderboard_entries(type, value DESC);

-- Per character leaderboard query  
CREATE INDEX IF NOT EXISTS idx_leaderboard_character ON leaderboard_entries(type, character, value DESC);

-- Player-specific leaderboard query
CREATE INDEX IF NOT EXISTS idx_leaderboard_player ON leaderboard_entries(player_id);

-- ============================================================================
-- Comments for documentation
-- ============================================================================
COMMENT ON TABLE players IS 'Hollow Wilds player accounts';
COMMENT ON COLUMN players.playfab_id IS 'Unique PlayFab identifier for the player';
COMMENT ON COLUMN players.last_seen_at IS 'Last activity timestamp for the player';

COMMENT ON TABLE player_saves IS 'Player game save data with version control';
COMMENT ON COLUMN player_saves.save_version IS 'Version number for optimistic locking';
COMMENT ON COLUMN player_saves.save_data IS 'Complete game state as JSONB (world, player, inventory, etc.)';

COMMENT ON TABLE player_save_backups IS 'Historical save backups for recovery';
COMMENT ON COLUMN player_save_backups.save_version IS 'Save version at time of backup';
COMMENT ON COLUMN player_save_backups.save_data IS 'Snapshot of game state at backup time';

COMMENT ON TABLE analytics_events IS 'Analytics event log for Hollow Wilds and legacy systems';
COMMENT ON COLUMN analytics_events.event_type IS 'Event type (player_death, item_crafted, boss_killed, etc.)';
COMMENT ON COLUMN analytics_events.event_properties IS 'Event-specific data as JSONB';
COMMENT ON COLUMN analytics_events.session_id IS 'Session identifier for grouping events';
