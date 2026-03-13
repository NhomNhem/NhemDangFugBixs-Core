## ADDED Requirements

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
