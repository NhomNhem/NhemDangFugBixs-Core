## 1. Database Schema

- [x] 1.1 Create `level_configs` table (level_id, config_json, created_at, updated_at)
- [x] 1.2 Create `talent_configs` table (talent_id, config_json, created_at, updated_at)
- [x] 1.3 Add seed data for existing levels and talents
- [x] 1.4 Update `leaderboard_entries` table if necessary to support dynamic types

## 2. Persistence Layer Refactor

- [x] 2.1 Update `PostgresLeaderboardRepository` to use dynamic metric types (replace `TODO_TYPE`)
- [x] 2.2 Implement `GetConfig` in `PostgresLevelRepository` using database queries
- [x] 2.3 Implement `GetConfig` in `PostgresTalentRepository` using database queries
- [x] 2.4 Implement short-term in-memory caching for configurations

## 3. Infrastructure & Middleware

- [x] 3.1 Implement Fiber custom error handler for standardized JSON responses
- [x] 3.2 Update `middleware/logger.go` to support JSON output in production
- [x] 3.3 Implement Redis-backed rate limiting middleware
- [x] 3.4 Apply rate limiting to high-volume endpoints (Analytics, Auth)

## 4. API Standardization

- [x] 4.1 Update all controllers to use the new standardized error response format
- [x] 4.2 Update Swagger annotations to reflect new error response schemas
- [x] 4.3 Update `docs/api.md` with standardized error documentation

## 5. Verification & Testing

- [x] 5.1 Add unit tests for database configuration retrieval
- [x] 5.2 Add integration tests for rate limiting (simulating burst traffic)
- [x] 5.3 Verify structured logging output in a production-like environment
- [x] 5.4 Update Unity client models to match the new error format (documented only)
