## ADDED Requirements

### Requirement: Verified Identity Assurance
The system SHALL use the PlayFab ID returned by the identity provider during ticket validation as the definitive source of truth for the player session.

#### Scenario: Identity binding
- **WHEN** a player logs in with a PlayFab ticket
- **THEN** the system ignores any client-provided IDs and binds the JWT strictly to the ID returned by the PlayFab API
