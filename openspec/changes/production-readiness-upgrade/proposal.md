## Why

The current backend implementation contains several "Phase 1" artifacts that rely on hardcoded logic, placeholder database queries, and lack essential production safeguards. To ensure stability, security, and scalability for actual gameplay, we must replace these placeholders with robust, database-driven implementations and implement production-standard middleware.

## What Changes

- Replace the hardcoded `TODO_TYPE` in the Leaderboard repository with dynamic metric handling.
- Transition `LevelRepository` and `TalentRepository` from hardcoded logic to actual PostgreSQL database queries.
- Implement Redis-backed rate limiting for high-volume endpoints (Analytics, Auth).
- Standardize API error responses to match the global snake_case format defined in the design document.
- **BREAKING**: Unified error response format may change existing error payload structures for some clients.

## Capabilities

### New Capabilities
- `data-validation`: Formal requirements for validating incoming save and analytics data.
- `production-ready`: Infrastructure requirements for rate limiting and logging.

### Modified Capabilities
- `leaderboard`: Update requirements to support multiple metric types dynamically.
- `level-progression`: Update to require persistent configuration storage.
- `talent-system`: Update to require persistent configuration storage.

## Impact

- **Affected code**: `internal/infrastructure/persistence/`, `internal/middleware/`, `internal/delivery/http/`.
- **APIs**: All endpoints under `/api/v1/` will adopt the new standardized error format.
- **Dependencies**: Redis client usage will expand to include rate-limiting keys.
- **Systems**: Supabase schema will require new tables or rows for Level and Talent configurations.
