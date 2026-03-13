## ADDED Requirements

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
