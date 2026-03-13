## ADDED Requirements

### Requirement: Get Global Leaderboard
The system SHALL return a paginated list of top players for a specific level ordered by completion time.

#### Scenario: Successful global leaderboard fetch
- **WHEN** player requests global leaderboard for a level
- **THEN** system returns top players ranked by best completion time with pagination support

#### Scenario: Invalid level ID
- **WHEN** player requests leaderboard for non-existent level
- **THEN** system returns 404 error with appropriate message

#### Scenario: Pagination beyond available data
- **WHEN** player requests page beyond available results
- **THEN** system returns empty array with total count metadata

### Requirement: Get Player Rank
The system SHALL return a specific player's rank and score for a given level.

#### Scenario: Player has completed the level
- **WHEN** player requests their rank for a completed level
- **THEN** system returns their rank, best time, and surrounding players (rank-1 and rank+1)

#### Scenario: Player has not completed the level
- **WHEN** player requests their rank for an uncompleted level
- **THEN** system returns null rank with level completion prompt

### Requirement: Filter Leaderboard by Time Period
The system SHALL support filtering leaderboards by time period (all-time, weekly, daily).

#### Scenario: Weekly leaderboard request
- **WHEN** player requests weekly leaderboard
- **THEN** system returns only completions from the past 7 days

#### Scenario: Daily leaderboard request
- **WHEN** player requests daily leaderboard
- **THEN** system returns only completions from the past 24 hours
