## ADDED Requirements

### Requirement: Post-Deployment Smoke Tests
The system SHALL support a suite of non-destructive verification scripts for use immediately after deployment.

#### Scenario: Production verification
- **WHEN** the deployment pipeline finishes a release to the `production` environment
- **THEN** it executes the `test_hollow_wilds.go` script against the live endpoint

### Requirement: System Connection Validation
The smoke tests SHALL verify connectivity to all primary data sources and external APIs.

#### Scenario: Connection testing
- **WHEN** the `test-db-connection.go` script is run in the production environment
- **THEN** it confirms that the backend can read from the `users` and `players` tables
