## MODIFIED Requirements

### Requirement: Character and Origin Selection
The system SHALL allow players to select one of 4 valid characters (RIMBA, DARA, BAYU, SARI) and a valid Origin (Hutan, Pantai, Gunung, Rawa) during the first login, strictly enforcing these choices against a whitelist.

#### Scenario: Initial setup with valid choices
- **WHEN** a new player authenticates and provides valid character/origin names
- **THEN** the system initializes their profile successfully

#### Scenario: Initial setup with invalid choices
- **WHEN** a new player provides an unknown character name
- **THEN** the system returns a `validation_error` (422)
