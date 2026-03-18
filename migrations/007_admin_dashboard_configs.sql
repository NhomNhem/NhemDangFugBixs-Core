-- level_configs: admin-managed level configurations
CREATE TABLE IF NOT EXISTS level_configs (
    level_id         TEXT PRIMARY KEY,
    map_id           TEXT NOT NULL,
    name             TEXT NOT NULL,
    difficulty       INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    min_time_seconds FLOAT NOT NULL,
    base_gold        INT NOT NULL DEFAULT 0,
    reward_stars     INT NOT NULL DEFAULT 3,
    objectives       JSONB NOT NULL DEFAULT '[]',
    is_active        BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- talent_configs: admin-managed talent configurations
CREATE TABLE IF NOT EXISTS talent_configs (
    talent_id       TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    max_level       INT NOT NULL,
    base_cost       INT NOT NULL,
    cost_scaling    FLOAT NOT NULL,
    bonus_per_level FLOAT NOT NULL,
    stat_type       TEXT NOT NULL,
    unlock_map      INT NOT NULL DEFAULT 1,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
