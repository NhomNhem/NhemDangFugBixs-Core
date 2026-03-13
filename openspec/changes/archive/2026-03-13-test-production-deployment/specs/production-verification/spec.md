## ADDED Requirements

### Requirement: Live Ticket Validation
The system SHALL successfully validate a real PlayFab session ticket against the live production environment.

#### Scenario: Successful production login
- **WHEN** client provides a valid real PlayFab session ticket to `https://gamefeel-backend.fly.dev/api/v1/auth/hw/login`
- **THEN** system returns a 200 OK with a valid HW JWT and the correct player ID

### Requirement: Live Data Persistence
The system SHALL correctly persist and retrieve data from the production Supabase database.

#### Scenario: Production save and load
- **WHEN** a player saves their state to the production endpoint
- **THEN** subsequent load requests return the exact state persisted in the production database
