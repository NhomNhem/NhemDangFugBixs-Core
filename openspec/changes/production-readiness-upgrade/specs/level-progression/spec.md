## MODIFIED Requirements

### Requirement: Level Configuration Retrieval
The system SHALL retrieve level-specific configurations (completion thresholds, base rewards) from the database.

#### Scenario: Getting level config
- **WHEN** a completion request is received
- **THEN** the system queries the `level_configs` table to fetch validation parameters for that level

### Requirement: Level Completion Persistence
The system SHALL record the results of a completed level, including stars earned and completion time.

#### Scenario: Submitting level results
- **WHEN** a player completes a level and submits the result
- **THEN** the system updates the `level_completions` table and returns the player's performance stats
