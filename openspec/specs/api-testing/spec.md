## ADDED Requirements

### Requirement: Automated Integration Tests
The system SHALL provide a suite of integration tests that verify the core Hollow Wilds Phase 1 endpoints.

#### Scenario: Successful test execution
- **WHEN** the test script in `scripts/test_hollow_wilds.go` is executed against a running server
- **THEN** all test cases (Auth, Save, Leaderboard, Analytics) pass successfully

### Requirement: End-to-End Flow Verification
The system SHALL verify that a player can log in, save their state, and submit a leaderboard entry in a single sequence.

#### Scenario: Complete player lifecycle
- **WHEN** a test player performs a full sequence of actions
- **THEN** the system maintains state correctly across all operations

### Requirement: API Integration Tests
The system SHALL provide a framework for testing HTTP endpoints by sending real requests to a test instance of the server.

#### Scenario: Testing health check endpoint
- **WHEN** the integration test sends a `GET /health` request
- **THEN** the system returns a `200 OK` status and a valid JSON response matching the health schema

### Requirement: Test Environment Configuration
The system SHALL allow running API tests against a configurable base URL and environment.

#### Scenario: Running tests against local server
- **WHEN** the test suite is configured with `BASE_URL=http://localhost:8080`
- **THEN** all API tests target the local development server
