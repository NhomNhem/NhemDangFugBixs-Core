## MODIFIED Requirements

### Requirement: Talent Retrieval
The system SHALL return the current talent levels for an authenticated player.

#### Scenario: Getting player talents
- **WHEN** a player requests their talents
- **THEN** the system returns a list of talent IDs and their corresponding current levels

### Requirement: Talent Upgrading
The system SHALL allow players to spend resources to upgrade their talents up to a maximum level.

#### Scenario: Successful talent upgrade
- **WHEN** a player has sufficient resources and submits an upgrade request
- **THEN** the system increments the talent level and persists the change to the `user_talents` table

## ADDED Requirements

### Requirement: Talent Max Level Enforcement
The system SHALL prevent players from upgrading talents beyond their defined `maxLevel`.

#### Scenario: Attempting to upgrade maxed talent
- **WHEN** a player submits an upgrade request for a talent already at `maxLevel`
- **THEN** the system returns a `validation_error` (422) and does not deduct any gold
