## ADDED Requirements

### Requirement: Request Correlation
The system SHALL generate and propagate a unique `request_id` for every incoming HTTP request.

#### Scenario: Correlating logs for a request
- **WHEN** a client makes an API request to the backend
- **THEN** all log entries related to that request include the same `X-Request-ID` value

### Requirement: JSON-Formatted Logging
The system SHALL support structured, machine-readable logs in JSON format for production environments.

#### Scenario: Production log ingestion
- **WHEN** the server is running in `production` mode
- **THEN** every log line is a valid JSON object containing `level`, `msg`, `time`, and `request_id`
