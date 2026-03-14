## ADDED Requirements

### Requirement: Global Rate Limiting
The system SHALL limit the number of requests per IP address to 100 requests per minute.

#### Scenario: Exceeding rate limit
- **WHEN** a client makes 101 requests within 60 seconds from the same IP
- **THEN** the system returns a `rate_limited` (429) error for subsequent requests until the window resets

### Requirement: Standardized Error Responses
The system SHALL return all errors in a unified JSON format containing `error` (code), `message`, and `trace_id`.

#### Scenario: Encountering an internal error
- **WHEN** an unhandled exception occurs
- **THEN** the system returns a `500 Internal Server Error` with a `snake_case` error code and a unique `trace_id`

### Requirement: Structured Logging
The system SHALL output logs in JSON format during production to facilitate log aggregation and analysis.

#### Scenario: Production log output
- **WHEN** the `ENV` variable is set to "production"
- **THEN** all logs are emitted as single-line JSON objects containing timestamps and severity levels
