## MODIFIED Requirements

### Requirement: Survival State Persistence
The system SHALL store and retrieve complete game states including Health, Hunger, Sanity, and Warmth for each player, ensuring all values are within DESIGNATED safety bounds.

#### Scenario: Successful save
- **WHEN** client sends a valid JSON save body with stats within bounds (0-100)
- **THEN** system persists the state via the `SaveRepository` and returns a new `save_version`

#### Scenario: Out of bounds save
- **WHEN** client sends a save body with Health > 100
- **THEN** system rejects the save with a `validation_error` (422)

## ADDED Requirements

### Requirement: Cache Expiration Policy
The system SHALL enforce a Time-To-Live (TTL) of 5 minutes for player save data cached in Redis.

#### Scenario: Automatic cache clearing
- **WHEN** 5 minutes have passed since the last cache update for a player
- **THEN** the record is automatically removed from Redis memory
