## MODIFIED Requirements

### Requirement: Level Completion Persistence
The system SHALL record the results of a completed level, including stars earned and completion time.

#### Scenario: Submitting level results
- **WHEN** a player completes a level and submits the result
- **THEN** the system updates the `level_completions` table and returns the player's performance stats

### Requirement: Statistics Aggregation
The system SHALL calculate and return global and player-specific statistics for level performance.

#### Scenario: Getting level stats
- **WHEN** a client requests statistics for a specific level ID
- **THEN** the system returns total plays, average time, and average stars

## ADDED Requirements

### Requirement: New Record Detection
The system SHALL determine if a level completion represents a new "personal best" for the player's time.

#### Scenario: Achieving new best time
- **WHEN** a player completes a level in 45 seconds and their previous best was 60 seconds
- **THEN** the `newBestTime` flag in the response is set to `true`
