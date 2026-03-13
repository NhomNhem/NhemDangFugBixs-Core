## ADDED Requirements

### Requirement: Service Health Checks
The system SHALL provide an `/health` endpoint for monitoring and readiness checks.

#### Scenario: Readiness check
- **WHEN** the orchestrator (Fly.io) checks the `/health` endpoint
- **THEN** the system returns a `200 OK` only if database and Redis connections are established

### Requirement: Error Tracking
The system SHALL log all unexpected errors with a unique `request_id` for traceability.

#### Scenario: Server error occurrence
- **WHEN** a handler catches an unexpected server error
- **THEN** it logs the error details, stack trace, and `request_id` before returning a `500` response
