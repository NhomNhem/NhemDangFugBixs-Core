## ADDED Requirements

### Requirement: Save Data Schema Validation
The system SHALL validate that player save data contains all required sections (world, player, inventory, sebilah, base).

#### Scenario: Submitting incomplete save data
- **WHEN** a player submits a save request missing the `inventory` section
- **THEN** the system returns a `validation_error` (422) and does not persist the data

### Requirement: Character Enum Validation
The system SHALL ensure that the player character name belongs to the approved list (RIMBA, DARA, BAYU, SARI).

#### Scenario: Using an invalid character
- **WHEN** a player submits save data with character name "UNKNWON"
- **THEN** the system returns a `validation_error` (422) with details about allowed character names

### Requirement: Analytics Event Payload Validation
The system SHALL validate analytics payloads against the expected properties for each event type.

#### Scenario: Malformed event payload
- **WHEN** a client submits a `player_death` event missing the `cause` property
- **THEN** the system rejects the event and increments the `rejected` counter in the response
