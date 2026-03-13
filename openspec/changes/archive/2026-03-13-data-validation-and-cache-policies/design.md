## Context

The system currently relies on the Unity client to send accurate game state data. While we trust the client for gameplay, the backend must verify that the state is within reasonable limits to prevent accidental corruption or intentional manipulation (cheating). Additionally, our Redis caching lacks expiration policies, leading to unnecessary memory usage.

## Goals / Non-Goals

**Goals:**
- Implement automated request body validation using tags.
- Enforce business-level constraints on survival stats.
- Define a global 5-minute TTL for cached player data.
- Standardize the error response for validation failures.

**Non-Goals:**
- Advanced anti-cheat (e.g., verifying move speed or physics).
- Complex multi-level caching strategies.

## Decisions

### 1. Validation Library: `go-playground/validator`
**Decision:** Use the standard `validator` package for Go.
**Rationale:** It is highly performant, widely adopted, and allows us to define rules using struct tags, which keeps the DTOs clean and readable.

### 2. Validation Layers: Two-Tier approach
**Decision:** 
- **Delivery Layer**: Validates syntax and basic bounds (e.g., "Health must be between 0 and 100").
- **Usecase Layer**: Validates stateful rules (e.g., "Player cannot upgrade Sebilah without required essence").
**Rationale:** Separation of concerns. Syntax/bounds checks belong at the edge, while business rules belong in the core logic.

### 3. Cache TTL: Global 5-minute expiration
**Decision:** Set a 300-second TTL on all `player:save:{id}` keys in Redis.
**Rationale:** Saves are loaded primarily at start or on cross-server transition (future). 5 minutes is enough to handle bursts without keeping memory locked indefinitely.

## Risks / Trade-offs

- **[Risk] Strict Validation Breaks Client** → **[Mitigation]** Ensure limits match the Unity client's configuration exactly and provide clear error messages so client logs can help debug.
- **[Trade-off] Performance Overhead** → Running validation on every save adds a few milliseconds. **[Mitigation]** This is negligible compared to database I/O and provides much higher safety.
