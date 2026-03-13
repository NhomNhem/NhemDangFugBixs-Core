## Leaderboard Specification

### Requirement: Submit Leaderboard Entry
The system SHALL create or update a player's leaderboard entry when they complete a level.

#### Scenario: New personal best time
- **WHEN** player submits a level completion with a time faster than their existing best
- **THEN** system updates their leaderboard entry with the new best time and recalculates rank

#### Scenario: Worse than personal best
- **WHEN** player submits a level completion slower than their existing best
- **THEN** system ignores the submission for leaderboard purposes (keeps existing best)

#### Scenario: First completion of a level
- **WHEN** player completes a level for the first time
- **THEN** system creates a new leaderboard entry with their completion time

### Requirement: Validate Submission Before Recording
The system SHALL validate that a level completion meets all validation rules, including anti-cheat checks, before recording to the leaderboard.

#### Scenario: Invalid completion data
- **WHEN** player submits level completion that fails validation (e.g., anti-cheat violation)
- **THEN** system rejects the submission and does not update leaderboard

#### Scenario: Duplicate submission same time
- **WHEN** player submits identical completion data twice
- **THEN** system treats as duplicate and does not create additional entry

### Requirement: Invalidate Cache on Update
The system SHALL invalidate cached leaderboard data when a player's entry changes.

#### Scenario: Top 1000 player updates score
- **WHEN** a player in the top 1000 updates their best time
- **THEN** system immediately invalidates the cached leaderboard for that level

#### Scenario: Player outside top 1000 updates score
- **WHEN** a player outside top 1000 updates their best time
- **THEN** system updates database but does not invalidate cache (cache will expire naturally)

### Requirement: Get Global Leaderboard
The system SHALL return a paginated list of top players for a specific level ordered by completion time.

#### Scenario: Successful global leaderboard fetch
- **WHEN** player requests global leaderboard for a level with `page` and `perPage` parameters
- **THEN** system returns top players ranked by best completion time with pagination support

#### Scenario: Invalid level ID
- **WHEN** player requests leaderboard for non-existent level
- **THEN** system returns 404 error with appropriate message

#### Scenario: Pagination beyond available data
- **WHEN** player requests page beyond available results
- **THEN** system returns empty array with total count metadata

### Requirement: Filter Leaderboard by Time Period
The system SHALL support filtering leaderboards by time period (all-time, weekly, daily).

#### Scenario: Weekly leaderboard request
- **WHEN** player requests weekly leaderboard
- **THEN** system returns only completions from the past 7 days

#### Scenario: Daily leaderboard request
- **WHEN** player requests daily leaderboard
- **THEN** system returns only completions from the past 24 hours

### Requirement: Get Player Rank
The system SHALL return a specific player's rank and score for a given level.

#### Scenario: Player has completed the level
- **WHEN** player requests their rank for a completed level
- **THEN** system returns their rank, best time, and surrounding players (rank-1 and rank+1)

#### Scenario: Player has not completed the level
- **WHEN** player requests their rank for an uncompleted level
- **THEN** system returns null rank with level completion prompt

### Requirement: Get Friends Leaderboard
The system SHALL return a leaderboard showing only the player's friends for a specific level using social graph integration.

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
