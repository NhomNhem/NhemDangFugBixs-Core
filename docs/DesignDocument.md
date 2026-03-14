# Hollow Wilds — Backend API Spec
**Version:** 1.0 · 2026-03-11  
**Server:** Go + Fiber + Supabase + Redis + Stripe  
**Base URL:** `https://api.hollowwilds.com/v1`  
**Auth:** PlayFab session ticket → JWT  
**Phase:** 1 (Auth, Save/Load, Analytics) — Multiplayer Phase 3

---

## Stack Context

```
Framework:  Fiber (Go)
Database:   Supabase (PostgreSQL)
Cache:      Upstash Redis
Auth:       PlayFab session ticket validation → JWT
Payments:   Stripe (Phase 2 cosmetics)
Client:     Unity 2022.3 + PlayFab SDK + UniTask
```

---

## Existing Endpoints (keep, không thay đổi)

```
GET  /health                          → health check
POST /api/v1/auth/login               → PlayFab auth → JWT
POST /api/v1/levels/complete          → level completion
POST /api/v1/talents/upgrade          → talent upgrade
POST /api/v1/payments/create-session  → Stripe checkout
POST /api/v1/analytics/events         → analytics
```

---

## New Endpoints — Phase 1

### 1. Auth

#### POST /api/v1/auth/login
*(existing — verify còn hoạt động)*

```json
Request:
{
  "playfab_session_ticket": "ABC123..."
}

Response 200:
{
  "token": "eyJhbGc...",
  "refresh_token": "ref_abc...",
  "expires_in": 3600,
  "player_id": "uuid-here"
}

Response 401:
{
  "error": "invalid_playfab_ticket",
  "message": "PlayFab session ticket invalid or expired"
}
```

#### POST /api/v1/auth/refresh

```json
Request:
{
  "refresh_token": "ref_abc..."
}

Response 200:
{
  "token": "eyJhbGc...",
  "expires_in": 3600
}

Response 401:
{
  "error": "invalid_refresh_token"
}
```

#### DELETE /api/v1/auth/logout
*(requires JWT header)*

```json
Response 200:
{
  "success": true
}
```

---

### 2. Player Save / Load

Game state được lưu per player. Deterministic world generation từ `world_seed` — server không lưu chunk data.

#### GET /api/v1/player/save
*(requires JWT)*

```json
Response 200:
{
  "player_id": "uuid",
  "save_version": 3,
  "updated_at": "2026-03-11T10:00:00Z",
  "world": {
    "seed": 1234567890,
    "play_time_seconds": 3600,
    "day_count": 5
  },
  "player": {
    "character": "RIMBA",
    "position": { "x": 128.5, "z": 64.2 },
    "health": 85.0,
    "hunger": 60.0,
    "sanity": 75.0,
    "warmth": 90.0
  },
  "inventory": {
    "slots": [
      { "slot": 0, "item_id": "wood_plank", "quantity": 12 },
      { "slot": 1, "item_id": "keris_lvl1", "quantity": 1 }
    ],
    "equipped_weapon": "keris_lvl1"
  },
  "sebilah": {
    "weapon_id": "keris_lvl1",
    "soul_level": 2,
    "infusion_points": 150
  },
  "base": {
    "placed_objects": [
      { "object_id": "campfire", "x": 100.0, "z": 50.0 }
    ]
  },
  "discovered_pois": ["shrine_01", "spirit_tree_03"],
  "quest_flags": {
    "tirai_fragment_1": true,
    "met_tok_batin": false
  }
}

Response 404:
{
  "error": "save_not_found",
  "message": "No save data found for this player"
}
```

#### PUT /api/v1/player/save
*(requires JWT)*

```json
Request: (same schema as GET response, minus player_id/updated_at/save_version)
{
  "world": { ... },
  "player": { ... },
  "inventory": { ... },
  "sebilah": { ... },
  "base": { ... },
  "discovered_pois": [...],
  "quest_flags": { ... }
}

Response 200:
{
  "success": true,
  "save_version": 4,
  "updated_at": "2026-03-11T10:05:00Z"
}

Response 409 (version conflict — multiplayer Phase 3):
{
  "error": "version_conflict",
  "server_version": 5,
  "message": "Save is outdated, fetch latest first"
}
```

#### POST /api/v1/player/save/backup
*(requires JWT — manual backup trigger)*

```json
Response 200:
{
  "success": true,
  "backup_id": "bkp_abc123",
  "created_at": "2026-03-11T10:05:00Z"
}
```

#### GET /api/v1/player/save/backups
*(requires JWT)*

```json
Response 200:
{
  "backups": [
    {
      "backup_id": "bkp_abc123",
      "save_version": 3,
      "created_at": "2026-03-11T10:05:00Z"
    }
  ]
}
```

---

### 3. Leaderboard

**3 types, locked for EA launch:**

| Type | Metric | Scope | Description |
|---|---|---|---|
| `longest_run_days` | int | global + per character | Ngày sống sót trong 1 run duy nhất |
| `sebilah_soul_level` | int (0-5) | global + per character | Soul level cao nhất đạt được |
| `bosses_killed` | int | global + per character | Tổng bosses killed across all runs |

> **Scope:** Mỗi type có 2 scopes — `global` (tất cả characters) và `per_character` (filter theo RIMBA/DARA/BAYU/SARI)

#### GET /api/v1/leaderboard
*(public, no auth required)*

```
Query params:
  type:      "longest_run_days" | "sebilah_soul_level" | "bosses_killed"
  scope:     "global" | "per_character" (default: global)
  character: "RIMBA" | "DARA" | "BAYU" | "SARI" (required khi scope=per_character)
  limit:     int (default 100, max 500)
  offset:    int (default 0)
```

```json
Response 200:
{
  "type": "longest_run_days",
  "scope": "per_character",
  "character": "RIMBA",
  "total": 350,
  "entries": [
    {
      "rank": 1,
      "player_id": "uuid",
      "display_name": "Player1",
      "value": 42,
      "character": "RIMBA",
      "world_seed": 1234567890,
      "combat_build": "berserker",
      "updated_at": "2026-03-11T10:00:00Z"
    }
  ]
}
```

> **Note:** `combat_build` và `world_seed` là metadata — players có thể verify/replay run của người khác.

#### POST /api/v1/leaderboard/submit
*(requires JWT — submit sau khi run kết thúc)*

```json
Request:
{
  "type": "longest_run_days",
  "value": 15,
  "character": "RIMBA",
  "world_seed": 1234567890,
  "combat_build": "balanced",
  "run_metadata": {
    "bosses_killed": 3,
    "biomes_explored": 4,
    "sebilah_soul_level": 2,
    "play_time_seconds": 7200
  }
}

Response 200:
{
  "success": true,
  "global_rank": 42,
  "character_rank": 15,
  "previous_global_rank": 55,
  "is_personal_best": true
}

Response 400:
{
  "error": "value_too_low",
  "message": "Submitted value does not beat personal best"
}
```

> **Anti-cheat note:** Server chỉ accept value nếu > personal best hiện tại. Không thể submit thấp hơn.

#### GET /api/v1/leaderboard/player
*(requires JWT — get own ranks across all types)*

```json
Response 200:
{
  "entries": [
    {
      "type": "longest_run_days",
      "global_rank": 42,
      "character_rank": 15,
      "character": "RIMBA",
      "value": 15,
      "personal_best": true
    },
    {
      "type": "sebilah_soul_level",
      "global_rank": 88,
      "character_rank": 30,
      "character": "RIMBA",
      "value": 3,
      "personal_best": true
    },
    {
      "type": "bosses_killed",
      "global_rank": 200,
      "character_rank": 67,
      "character": "RIMBA",
      "value": 8,
      "personal_best": true
    }
  ]
}

---

### 4. Analytics

#### POST /api/v1/analytics/events
*(existing — extend schema)*

```json
Request:
{
  "events": [
    {
      "event_name": "player_death",
      "timestamp": "2026-03-11T10:00:00Z",
      "session_id": "sess_abc",
      "payload": {
        "cause": "spirit_attack",
        "day_count": 5,
        "biome": "night_marsh",
        "character": "RIMBA",
        "build": "balanced"
      }
    },
    {
      "event_name": "item_crafted",
      "timestamp": "2026-03-11T10:01:00Z",
      "session_id": "sess_abc",
      "payload": {
        "item_id": "keris_lvl1",
        "materials_used": ["iron_shard", "spirit_wood"]
      }
    }
  ]
}

Response 200:
{
  "accepted": 2,
  "rejected": 0
}
```

**Standard event names:**
```
session_start       session_end
player_death        player_spawn
item_crafted        item_used
enemy_killed        boss_killed
biome_entered       poi_discovered
sebilah_evolved     build_changed
base_object_placed  quest_flag_set
```

---

### 5. Economy / Shop (Phase 2 — design now, implement later)

#### GET /api/v1/shop/items
```json
Response 200:
{
  "items": [
    {
      "item_id": "skin_rimba_sea_warrior",
      "name": "RIMBA — Sea Warrior",
      "price_usd": 4.99,
      "type": "character_skin",
      "preview_url": "https://cdn.hollowwilds.com/skins/..."
    }
  ]
}
```

#### POST /api/v1/payments/create-session
*(existing Stripe — extend for cosmetics)*

```json
Request:
{
  "item_id": "skin_rimba_sea_warrior",
  "success_url": "hollowwilds://purchase/success",
  "cancel_url": "hollowwilds://purchase/cancel"
}
```

---

## Supabase Schema

```sql
-- Players
CREATE TABLE players (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  playfab_id VARCHAR(64) UNIQUE NOT NULL,
  display_name VARCHAR(64),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  last_seen_at TIMESTAMPTZ DEFAULT NOW()
);

-- Save Data
CREATE TABLE player_saves (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id UUID REFERENCES players(id) ON DELETE CASCADE,
  save_version INT DEFAULT 1,
  save_data JSONB NOT NULL,          -- full game state JSON
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(player_id)
);

-- Save Backups
CREATE TABLE player_save_backups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id UUID REFERENCES players(id),
  save_version INT,
  save_data JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Leaderboard
-- type: longest_run_days | sebilah_soul_level | bosses_killed
-- Stores personal best per (player, type, character)
CREATE TABLE leaderboard_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id UUID REFERENCES players(id),
  type VARCHAR(32) NOT NULL,
  value BIGINT NOT NULL,
  character VARCHAR(16) NOT NULL,    -- RIMBA | DARA | BAYU | SARI
  world_seed BIGINT,
  combat_build VARCHAR(16),          -- balanced | berserker | shade_walker
  run_metadata JSONB,                -- bosses_killed, biomes_explored, etc
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(player_id, type, character) -- 1 personal best per player+type+character
);
-- Global leaderboard query
CREATE INDEX idx_leaderboard_global ON leaderboard_entries(type, value DESC);
-- Per character leaderboard query  
CREATE INDEX idx_leaderboard_character ON leaderboard_entries(type, character, value DESC);

-- Analytics
CREATE TABLE analytics_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id UUID REFERENCES players(id),
  session_id VARCHAR(64),
  event_name VARCHAR(64) NOT NULL,
  payload JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_analytics_event_name ON analytics_events(event_name, created_at);

-- Owned Items (Phase 2)
CREATE TABLE player_owned_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id UUID REFERENCES players(id),
  item_id VARCHAR(64) NOT NULL,
  purchased_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(player_id, item_id)
);
```

---

## Redis Cache Strategy

```
Key patterns:
  player:save:{player_id}          TTL: 300s  (save data cache)
  leaderboard:{type}:top100        TTL: 60s   (leaderboard cache)
  session:{jwt_jti}:blacklist      TTL: 3600s (logout blacklist)
  ratelimit:{ip}:{endpoint}        TTL: 60s   (rate limiting)
```

---

## Middleware Stack

```
All routes:
  RateLimit     → 100 req/min per IP
  RequestID     → trace ID per request
  Logger        → structured logging
  Recover       → panic recovery

Protected routes (JWT required):
  JWTAuth       → validate + extract player_id
  ActivityLog   → update last_seen_at (async)
```

---

## Error Response Format

All API errors follow a standardized JSON format:

```json
{
  "success": false,
  "error": {
    "code": "snake_case_error_code",
    "message": "Human readable message",
    "trace_id": "req_abc123"
  }
}
```

**Standard error codes:**
```
invalid_request       → 400 (bad request body/params)
unauthorized          → 401 (invalid/expired JWT)
forbidden             → 403 (valid JWT, wrong resource)
not_found             → 404
conflict              → 409 (version mismatch)
validation_error      → 422 (schema validation failed)
rate_limited          → 429
internal_error        → 500
```

---

## Unity Client Integration

```csharp
// BackendService.cs pattern (existing)
// Agent KHÔNG thay đổi client code
// Chỉ implement server endpoints

// Client sẽ call:
await _backendService.SaveGameAsync(saveData);
await _backendService.LoadGameAsync();
await _backendService.SubmitLeaderboardAsync(type, value);
await _backendService.TrackEventAsync(eventName, payload);
```

---

## Phase 3 Multiplayer — Design Notes (không implement Phase 1)

```
Khi Fish-Net integrate Phase 3:
  WebSocket endpoint: /ws/game/{session_id}
  Server authoritative cho: entity positions, combat
  Client authoritative cho: chunk generation (deterministic seed)
  
Save system extend:
  PUT /api/v1/player/save → add "session_id" field
  Conflict resolution: last-write-wins per field
  
Design server code extensible cho WebSocket:
  Tách HTTP handlers khỏi business logic
  Services layer không biết về transport layer
```

---

## Implementation Priority

```
Sprint 1 (tuần này):
  ✅ Auth refresh endpoint
  ✅ Save/Load endpoints
  ✅ Supabase schema migration

Sprint 2:
  ✅ Leaderboard endpoints
  ✅ Analytics extend schema
  ✅ Redis caching layer

Sprint 3 (Phase 2 prep):
  ✅ Shop/Economy endpoints
  ✅ Stripe webhook for cosmetics
  ✅ player_owned_items table
```

---

*Hollow Wilds · Backend API Spec v1.0 · 2026-03-11*