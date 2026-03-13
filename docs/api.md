# Hollow Wilds Backend API Documentation

**Base URL**: `http://localhost:8080/api/v1` (development)  
**Production URL**: `https://api.yourgame.com/api/v1`

## Authentication

Most endpoints require JWT authentication.

**Header:**
```
Authorization: Bearer <jwt_token>
```

Get JWT token via `/auth/hw/login` endpoint for Hollow Wilds.

---

## Endpoints

### Health Check

**GET /health**

Check if server is running and database/Redis connections are healthy.

---

### Authentication (Hollow Wilds)

**POST /api/v1/auth/hw/login**

Authenticate with PlayFab session ticket.

**Request Body:**
```json
{
  "playfab_session_ticket": "SESSION_TICKET_FROM_PLAYFAB"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1Ni...",
  "refresh_token": "uuid-v4-token",
  "expires_in": 3600,
  "player_id": "uuid-v4-player-id"
}
```

**POST /api/v1/auth/refresh**

Get new access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "uuid-v4-token"
}
```

**DELETE /api/v1/auth/logout**

Invalidate tokens.

---

### Player Save System

**GET /api/v1/player/save**

Load complete game state.

**Authentication:** Required

**Response:**
```json
{
  "player_id": "uuid",
  "save_version": 15,
  "updated_at": "2026-03-14T01:30:00Z",
  "world": { "seed": 123, "play_time_seconds": 3600, "day_count": 5 },
  "player": { "character": "RIMBA", "health": 85.0, ... },
  "inventory": { "slots": [...], "equipped_weapon": "sword_01" },
  "sebilah": { ... },
  "base": { "placed_objects": [...] },
  "discovered_pois": ["forest_shrine", "cave_entrance"],
  "quest_flags": { "tutorial_complete": true }
}
```

**PUT /api/v1/player/save?version={expected_version}**

Save game state with optimistic locking.

**Authentication:** Required

**Request Body:** (Same structure as Load response)

**Error (409 Conflict):**
```json
{
  "error": "version_conflict",
  "server_version": 15,
  "message": "Save is outdated, fetch latest first"
}
```

---

### Backups

**POST /api/v1/player/save/backup**

Manually create a save backup.

**GET /api/v1/player/save/backups**

List all backups for the player.

**POST /api/v1/player/save/restore**

Restore save from a backup ID.

**Request Body:**
```json
{
  "backup_id": "uuid-of-backup"
}
```

---

### Leaderboards

**GET /api/v1/leaderboard**

Get Hollow Wilds ranked entries.

**Query Parameters:**
- `type`: Metric type (`longest_run_days`, `sebilah_soul_level`, `bosses_killed`). Default: `longest_run_days`.
- `scope`: Scope (`global`, `per_character`). Default: `global`.
- `character`: Character name (required if scope=per_character).
- `limit`: Number of entries. Default: 100.
- `offset`: Offset for pagination. Default: 0.

**GET /api/v1/leaderboards/{levelId}**

Get global rankings for a specific level.

**Query Parameters:**
- `page`: Page number. Default: 1.
- `perPage`: Entries per page. Default: 10.

**GET /api/v1/leaderboards/{levelId}/me**

Get authenticated player's rank and surrounding players for a level.

**Authentication:** Required

**GET /api/v1/leaderboards/{levelId}/friends**

Get friends rankings for a specific level.

**Authentication:** Required

---

### Admin Operations

All admin endpoints require a JWT token with admin privileges.

**GET /api/v1/admin/users/search?q={query}**

Search for users.

**GET /api/v1/admin/users/{userId}/profile**

Get detailed user profile.

**POST /api/v1/admin/users/{userId}/adjust-gold**

Adjust user gold balance.

**DELETE /api/v1/admin/leaderboards/{levelId}**

Reset all entries for a specific level leaderboard.

**Request Body:**
```json
{
  "reason": "Reason for reset (min 10 chars)"
}
```

**GET /api/v1/admin/leaderboards/stats**

Get analytics and statistics for all level leaderboards.

