## ADDED Requirements

### Requirement: Centralized Request Validation
The system SHALL provide a centralized mechanism to validate all inbound HTTP request bodies against predefined rules.

#### Scenario: Malformed JSON
- **WHEN** client sends a request with invalid JSON syntax
- **THEN** system returns a `validation_error` (422) with a descriptive message

### Requirement: Survival Stat Bounds Enforcement
The system SHALL enforce strict numerical bounds on all player survival metrics (Health, Hunger, Sanity, Warmth).

#### Scenario: Impossible health value
- **WHEN** client attempts to save a health value greater than 100 or less than 0
- **THEN** system rejects the request with a `validation_error` (422)

### Requirement: Character and Origin Whitelisting
The system SHALL ensure that only valid character names and origins defined in the game design are accepted.

#### Scenario: Invalid character selection
- **WHEN** a player attempts to select a character not in (RIMBA, DARA, BAYU, SARI)
- **THEN** system returns a `validation_error` (422)
