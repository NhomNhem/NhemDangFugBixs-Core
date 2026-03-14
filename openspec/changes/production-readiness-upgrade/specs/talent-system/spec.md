## MODIFIED Requirements

### Requirement: Talent Definition Retrieval
The system SHALL retrieve talent definitions (max levels, upgrade costs) from the database.

#### Scenario: Getting talent config
- **WHEN** an upgrade request is received
- **THEN** the system queries the `talent_configs` table to verify upgrade feasibility and costs

### Requirement: Talent Retrieval
The system SHALL return the current talent levels for an authenticated player.

#### Scenario: Getting player talents
- **WHEN** a player requests their talents
- **THEN** the system returns a list of talent IDs and their corresponding current levels
