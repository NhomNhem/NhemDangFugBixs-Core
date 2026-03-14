-- Migration: Game Configurations
-- Description: Adds tables for level and talent configurations
-- Date: 2026-03-14

-- ============================================================================
-- 1. Create level_configs table
-- ============================================================================
CREATE TABLE IF NOT EXISTS level_configs (
    level_id TEXT PRIMARY KEY,
    map_id TEXT,
    config_json JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================================
-- 2. Create talent_configs table
-- ============================================================================
CREATE TABLE IF NOT EXISTS talent_configs (
    talent_id TEXT PRIMARY KEY,
    config_json JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================================
-- 3. Seed data for existing levels
-- ============================================================================
INSERT INTO level_configs (level_id, map_id, config_json)
VALUES 
('level_1', 'map_1', '{
    "min_time_seconds": 10.0,
    "base_gold": 100,
    "objectives": [
        {"type": "completion", "threshold": 1, "operator": "gte"},
        {"type": "health", "threshold": 50, "operator": "gte"},
        {"type": "time", "threshold": 60, "operator": "lte"}
    ]
}'),
('level_2', 'map_1', '{
    "min_time_seconds": 15.0,
    "base_gold": 150,
    "objectives": [
        {"type": "completion", "threshold": 1, "operator": "gte"},
        {"type": "health", "threshold": 60, "operator": "gte"},
        {"type": "time", "threshold": 90, "operator": "lte"}
    ]
}')
ON CONFLICT (level_id) DO UPDATE SET config_json = EXCLUDED.config_json, updated_at = NOW();

-- ============================================================================
-- 4. Seed data for existing talents
-- ============================================================================
INSERT INTO talent_configs (talent_id, config_json)
VALUES 
('strength', '{"max_level": 5, "base_cost": 100, "cost_multiplier": 1.5}'),
('agility', '{"max_level": 5, "base_cost": 100, "cost_multiplier": 1.5}'),
('intelligence', '{"max_level": 5, "base_cost": 100, "cost_multiplier": 1.5}'),
('stamina', '{"max_level": 10, "base_cost": 50, "cost_multiplier": 1.2}')
ON CONFLICT (talent_id) DO UPDATE SET config_json = EXCLUDED.config_json, updated_at = NOW();
