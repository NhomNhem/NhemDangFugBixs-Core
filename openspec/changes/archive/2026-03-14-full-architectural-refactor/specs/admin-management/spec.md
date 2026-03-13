## ADDED Requirements

### Requirement: Admin Authentication and Authorization
The system SHALL require a valid JWT with `is_admin` claims for all administrative endpoints.

#### Scenario: Unauthorized admin access
- **WHEN** a non-admin user attempts to access `/api/v1/admin/users/search`
- **THEN** the system returns a `forbidden` (403) error

### Requirement: User Management
The system SHALL provide capabilities to search, profile, and adjust balances for any player.

#### Scenario: Adjusting player gold
- **WHEN** an admin submits a request to adjust a player's gold
- **THEN** the system updates the `users` (legacy) or `players` (Hollow Wilds) table and records the action

### Requirement: Ban Enforcement
The system SHALL allow admins to ban and unban users, preventing them from accessing the game.

#### Scenario: Banning a user
- **WHEN** an admin bans a user with a specific reason
- **THEN** the user's `is_banned` flag is set to true and they receive a `forbidden` error on next login
