## Context

The backend is transitioning from a functional prototype to a production-ready service. The current Phase 1 implementation served its purpose for initial integration but lacks the necessary flexibility and safeguards for a live game environment. Specifically, the system relies on hardcoded configurations for levels and talents, which makes tuning and adding content difficult. Additionally, the lack of centralized rate limiting and standardized error handling poses operational risks.

## Goals / Non-Goals

**Goals:**
- Eliminate all hardcoded "Phase 1" placeholders in the persistence layer.
- Centralize level and talent configurations in PostgreSQL.
- Implement Redis-backed global rate limiting.
- Standardize the error response format across all handlers.
- Enable JSON structured logging for production environments.

**Non-Goals:**
- Implementing a full Admin Dashboard (handled by direct Supabase access for now).
- Changing the core authentication flow (PlayFab → JWT).
- Performance optimization beyond basic indexing and caching already planned.

## Decisions

- **Configuration Storage**: Use PostgreSQL tables (`level_configs`, `talent_configs`) for all game balance data.
  - *Rationale*: Allows designers to tune the game via Supabase UI without redeploying code.
- **Rate Limiting**: Use the Fiber `limiter` middleware with a Redis storage backend.
  - *Rationale*: Redis provides a shared state across multiple server instances (scalability) and is already part of our stack.
- **Error Handling**: Implement a custom Fiber Error Handler that wraps all responses in the `APIResponse` struct.
  - *Rationale*: Ensures consistency and automatically includes `trace_id` for debugging.
- **Metric Types**: Update `LeaderboardRepository` to accept a `MetricType` string instead of hardcoded constants.
  - *Rationale*: Future-proofs the system for new types of leaderboards.

## Risks / Trade-offs

- **Risk**: Moving configs to the DB adds a query per request.
  - *Mitigation*: Implement short-term (e.g., 5 min) caching in Redis or application memory for configuration data.
- **Risk**: Breaking changes in error format for existing Unity clients.
  - *Mitigation*: Communicate the change clearly; update the Unity SDK models simultaneously.
