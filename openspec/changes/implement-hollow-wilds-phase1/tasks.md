## 1. Database Schema

- [x] 1.1 Create `players` table migration (id, playfab_id, display_name, created_at, last_seen_at)
- [x] 1.2 Create `player_saves` table migration (id, player_id, save_version, save_data JSONB, updated_at)
- [x] 1.3 Create `player_save_backups` table migration (id, player_id, save_version, save_data JSONB, created_at)
- [x] 1.4 Create indexes on players (playfab_id unique index)
- [x] 1.5 Create indexes on player_saves (player_id unique index)
- [x] 1.6 Create indexes on player_save_backups (player_id, created_at)
- [x] 1.7 Run migrations and verify schema in Supabase

## 2. Models

- [x] 2.1 Create `Player` model in internal/models/player.go
- [x] 2.2 Create `PlayerSave` model with JSONB save data structure
- [x] 2.3 Create `PlayerSaveBackup` model
- [x] 2.4 Create request DTOs: LoginRequest, RefreshTokenRequest, SaveGameRequest
- [x] 2.5 Create response DTOs: AuthResponse, SaveResponse, BackupResponse
- [x] 2.6 Define game state structs: World, PlayerState, Inventory, Sebilah, Base

## 3. Auth Service

- [x] 3.1 Implement PlayFab session ticket validation function
- [x] 3.2 Implement JWT generation with player_id claim
- [x] 3.3 Implement refresh token generation (UUID v4)
- [x] 3.4 Implement Redis storage for refresh tokens (7-day TTL)
- [x] 3.5 Implement refresh token validation function
- [x] 3.6 Implement token revocation on logout

## 4. Save Service

- [x] 4.1 Implement GetPlayerSave function with Redis caching
- [x] 4.2 Implement SavePlayerSave function with version control
- [x] 4.3 Implement version conflict detection
- [x] 4.4 Implement save data validation (required fields, character enum)
- [x] 4.5 Implement Redis cache invalidation on save
- [x] 4.6 Implement last_seen_at update on save operations

## 5. Backup Service

- [x] 5.1 Implement CreateBackup function
- [x] 5.2 Implement GetBackups function (list all backups for player)
- [x] 5.3 Implement RestoreFromBackup function
- [x] 5.4 Implement backup limit enforcement (max 10 backups)
- [x] 5.5 Implement automatic backup on major version thresholds
- [x] 5.6 Implement oldest backup deletion when limit reached

## 6. Auth Endpoints

- [x] 6.1 Update POST /api/v1/auth/login to accept playfab_session_ticket parameter
- [x] 6.2 Implement POST /api/v1/auth/refresh endpoint
- [x] 6.3 Implement DELETE /api/v1/auth/logout endpoint
- [x] 6.4 Add JWT middleware to protected routes
- [x] 6.5 Register new auth routes in internal/api/routes.go

## 7. Save Endpoints

- [x] 7.1 Implement GET /api/v1/player/save endpoint
- [x] 7.2 Implement PUT /api/v1/player/save endpoint
- [x] 7.3 Add version conflict error response (409)
- [x] 7.4 Add save not found error response (404)
- [x] 7.5 Register save routes in internal/api/routes.go

## 8. Backup Endpoints

- [x] 8.1 Implement POST /api/v1/player/save/backup endpoint
- [x] 8.2 Implement GET /api/v1/player/save/backups endpoint
- [x] 8.3 Implement POST /api/v1/player/save/restore endpoint
- [x] 8.4 Register backup routes in internal/api/routes.go

## 9. Analytics Extension

- [x] 9.1 Extend analytics event schema to support survival event types
- [x] 9.2 Update POST /api/v1/analytics/events to accept new event payloads
- [x] 9.3 Add rate limiting for analytics submissions (100 events/minute)
- [x] 9.4 Add session_id tracking for analytics events

## 10. Redis Integration

- [x] 10.1 Add Redis client initialization in pkg/utils/redis.go
- [x] 10.2 Implement save data cache (key: player:save:{player_id}, TTL: 300s)
- [x] 10.3 Implement refresh token storage (key: session:{token}, TTL: 7 days)
- [x] 10.4 Implement session blacklist for logout (key: session:{jti}:blacklist)

## 11. Testing

- [x] 11.1 Write unit tests for auth service (login, refresh, logout)
- [x] 11.2 Write unit tests for save service (save, load, version conflict)
- [x] 11.3 Write unit tests for backup service (create, list, restore)
- [x] 11.4 Write integration tests for auth endpoints
- [x] 11.5 Write integration tests for save endpoints
- [x] 11.6 Write integration tests for backup endpoints
- [x] 11.7 Test Redis caching behavior (cache hit/miss, invalidation)

## 12. Documentation

- [x] 12.1 Add Swagger annotations for new endpoints
- [x] 12.2 Update docs/api.md with Hollow Wilds endpoints
- [x] 12.3 Run `swag init` to regenerate swagger docs
- [x] 12.4 Add example requests/responses for save data structure
- [x] 12.5 Document error codes and response formats
