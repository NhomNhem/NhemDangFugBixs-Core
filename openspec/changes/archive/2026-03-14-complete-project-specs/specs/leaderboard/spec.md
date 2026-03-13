## ADDED Requirements

### Requirement: Global Rankings
The system SHALL maintain real-time leaderboards for all players based on level completion times and total stars earned.

#### Scenario: Submitting a new record
- **WHEN** a player completes a level with a faster time than their previous best
- **THEN** their position on the global leaderboard for that level is updated immediately

### Requirement: Friend-Based Leaderboards
The system SHALL provide a filtered view of the leaderboard containing only the player's friends.

#### Scenario: Requesting social leaderboard
- **WHEN** an authenticated user requests a leaderboard for a specific level with the `friends_only` filter
- **THEN** the system returns rankings only for players with matching `PlayFab` friend IDs
