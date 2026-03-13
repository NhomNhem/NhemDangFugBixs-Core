## ADDED Requirements

### Requirement: Standardized Metrics
The system SHALL track three core metrics for the Early Access launch: `longest_run_days`, `sebilah_soul_level`, and `bosses_killed`.

#### Scenario: Valid leaderboard query
- **WHEN** client requests a leaderboard with a supported `type`
- **THEN** system returns ranked entries from the repository matching that metric

### Requirement: Personal Best Verification
The system SHALL only update a player's entry if the submitted value exceeds their stored personal best for that specific character.

#### Scenario: Submitting a better run
- **WHEN** a player completes a run with a higher `days_survived` than their previous record for RIMBA
- **THEN** the system updates the `LeaderboardRepository` and returns the new ranks

### Requirement: Character and Scope Filtering
The system SHALL support filtering the leaderboard by character (RIMBA, DARA, BAYU, SARI) and scope (Global vs. Character-only).

#### Scenario: Character specific rank
- **WHEN** a player requests their rank for BAYU on the `longest_run_days` metric
- **THEN** the system returns both their global rank and their rank relative to other BAYU players
