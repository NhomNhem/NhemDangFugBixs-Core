## ADDED Requirements

### Requirement: Get Friends Leaderboard
The system SHALL return a leaderboard showing only the player's friends for a specific level.

#### Scenario: Successful friends leaderboard fetch
- **WHEN** player requests friends leaderboard for a level
- **THEN** system returns ranked list of only their friends who completed the level

#### Scenario: No friends have completed the level
- **WHEN** player requests friends leaderboard but no friends have completed the level
- **THEN** system returns empty array with message encouraging player to be first

#### Scenario: Player has no friends
- **WHEN** player with no friends requests friends leaderboard
- **THEN** system returns empty array with prompt to add friends

### Requirement: Include Player in Friends Leaderboard
The system SHALL always include the requesting player's rank in the friends leaderboard response.

#### Scenario: Player ranked among friends
- **WHEN** player requests friends leaderboard
- **THEN** system includes player's entry highlighted in the friends list

#### Scenario: Player not yet completed but friends have
- **WHEN** player hasn't completed level but friends have
- **THEN** system shows friends rankings with indicator that player hasn't completed

### Requirement: Refresh Friends List Cache
The system SHALL update the cached friends list when social connections change.

#### Scenario: Player adds a friend
- **WHEN** player adds a new friend
- **THEN** system invalidates cached friends leaderboard for all levels

#### Scenario: Friend removes player
- **WHEN** a friend removes the player from their friends list
- **THEN** system updates cached friends list on next leaderboard request
