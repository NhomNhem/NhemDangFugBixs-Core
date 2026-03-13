## ADDED Requirements

### Requirement: Batch Event Ingestion
The system SHALL support the ingestion of multiple analytics events in a single HTTP request to optimize client-server communication.

#### Scenario: Submitting game session events
- **WHEN** a player completes a 20-minute session and the client sends a batch of events
- **THEN** the system persists them via the `AnalyticsRepository` and returns the accepted count

### Requirement: Standardized Survival Events
The system SHALL recognize and validate a specific roster of survival-themed events (e.g., `player_death`, `item_crafted`, `enemy_killed`, `boss_killed`, `sebilah_evolved`).

#### Scenario: Tracking death
- **WHEN** a `player_death` event is received
- **THEN** the system logs the `cause`, `day_count`, and `character` for later heat-map or balance analysis

### Requirement: Event Context Enrichment
The system SHALL automatically attach the `player_id`, `session_id`, and `server_timestamp` to every incoming event.

#### Scenario: Searchable audit trail
- **WHEN** a developer queries the `analytics_events` table
- **THEN** every event can be correlated to a specific player and server-time for debugging
