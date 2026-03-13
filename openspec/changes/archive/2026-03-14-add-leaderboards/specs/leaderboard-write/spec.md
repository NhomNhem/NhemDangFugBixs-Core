## ADDED Requirements

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

### Requirement: Invalidate Cache on Update
The system SHALL invalidate cached leaderboard data when a player's entry changes.

#### Scenario: Top 1000 player updates score
- **WHEN** a player in the top 1000 updates their best time
- **THEN** system immediately invalidates the cached leaderboard for that level

#### Scenario: Player outside top 1000 updates score
- **WHEN** a player outside top 1000 updates their best time
- **THEN** system updates database but does not invalidate cache (cache will expire naturally)

### Requirement: Validate Submission Before Recording
The system SHALL validate that a level completion meets all validation rules before recording to leaderboard.

#### Scenario: Invalid completion data
- **WHEN** player submits level completion that fails validation
- **THEN** system rejects the submission and does not update leaderboard

#### Scenario: Duplicate submission same time
- **WHEN** player submits identical completion data twice
- **THEN** system treats as duplicate and does not create additional entry
