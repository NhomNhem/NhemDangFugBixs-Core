-- Legacy tables from GameFeel Generic Backend
-- These are included for backward compatibility

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
