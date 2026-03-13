## Context

The GameFeel backend currently tracks level completions but lacks competitive features. Players cannot compare their performance with others, reducing engagement and replayability. The system uses Fiber (Go), Supabase (PostgreSQL), and Upstash Redis for caching.

**Stakeholders:**
- Players: Want to compete and track progress
- Game designers: Need engagement metrics
- Admin team: Need leaderboard management tools

**Constraints:**
- Must integrate with existing JWT authentication
- Must work with current level completion flow
- Low-latency reads required for game UI responsiveness
- Must handle high write volume during peak play times

## Goals / Non-Goals

**Goals:**
- Implement global leaderboards for level completion times
- Support friend-based leaderboards using social connections
- Cache leaderboard queries for sub-100ms response times
- Provide admin tools for leaderboard management
- Maintain data integrity with existing level validation

**Non-Goals:**
- Real-time leaderboard updates (eventual consistency acceptable)
- Cross-platform leaderboards (Unity-only for now)
- Seasonal or reset-based leaderboards (future enhancement)
- Spectator features or replay sharing

## Decisions

### Database Schema Design
**Decision:** Use separate `leaderboard_entries` table with composite indexes

**Rationale:** 
- Keeps leaderboard data separate from core player data
- Enables efficient range queries for ranking
- Supports multiple leaderboard types (global, friends, weekly)

**Alternatives considered:**
- Materialized views: Too complex for incremental updates
- Redis-only storage: Risk of data loss, need PostgreSQL as source of truth

### Caching Strategy
**Decision:** Cache top 1000 entries per leaderboard with 30-second TTL

**Rationale:**
- Most players only care about top rankings and their own position
- 30-second TTL balances freshness with performance
- Redis sorted sets enable efficient rank queries

**Alternatives considered:**
- No caching: Too slow for real-time game UI
- Cache everything: Wastes memory on rarely-accessed data

### Leaderboard Update Timing
**Decision:** Update leaderboards synchronously after level completion validation

**Rationale:**
- Ensures player sees their score immediately
- Simplifies error handling (no eventual consistency confusion)
- Acceptable latency for write path

**Alternatives considered:**
- Async queue: Adds complexity, player might not see score
- Batch updates: Too much delay for player feedback

## Risks / Trade-offs

**Risk:** High write volume during peak times could stress database
→ **Mitigation:** Use Redis as write buffer with periodic flush to PostgreSQL

**Risk:** Friend leaderboards require social graph queries
→ **Mitigation:** Cache friend lists separately, update on social changes

**Risk:** Cheating through manipulated submissions
→ **Mitigation:** Rely on existing level validation; add anomaly detection later

**Trade-off:** Synchronous updates add latency to level completion
→ **Acceptable:** Players expect delay; correctness over speed

**Trade-off:** Caching means stale data for ~30 seconds
→ **Acceptable:** Leaderboards don't need real-time accuracy
