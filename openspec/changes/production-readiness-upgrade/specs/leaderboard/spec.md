## MODIFIED Requirements

### Requirement: Submit Leaderboard Entry
The system SHALL create or update a player's leaderboard entry when they complete a level or run.

#### Scenario: Submitting with dynamic metric type
- **WHEN** a player submits a result with a valid metric type (`longest_run_days`, `sebilah_soul_level`, `bosses_killed`)
- **THEN** the system persists the value to the corresponding metric category in the database

#### Scenario: New personal best time
- **WHEN** player submits a level completion with a time faster than their existing best
- **THEN** system updates their leaderboard entry with the new best time and recalculates rank

#### Scenario: Worse than personal best
- **WHEN** player submits a level completion slower than their existing best
- **THEN** system ignores the submission for leaderboard purposes (keeps existing best)

#### Scenario: First completion of a level
- **WHEN** player completes a level for the first time
- **THEN** system creates a new leaderboard entry with their completion time
