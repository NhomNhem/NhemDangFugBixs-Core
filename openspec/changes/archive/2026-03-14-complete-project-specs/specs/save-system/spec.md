## ADDED Requirements

### Requirement: Save State Persistence
The system SHALL allow players to save and retrieve their game progress as a serialized JSON blob in the `player_saves` table.

#### Scenario: Upserting save data
- **WHEN** an authenticated player submits a new save state
- **THEN** the system updates the existing record or creates a new one for that player

### Requirement: Automated Backups
The system SHALL automatically create a backup of the current save state before performing an update.

#### Scenario: Updating save with backup
- **WHEN** a player updates their save state
- **THEN** the current state is copied to the `player_save_backups` table before being overwritten
