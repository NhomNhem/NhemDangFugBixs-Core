## ADDED Requirements

### Requirement: Standardized Log Format
The system SHALL output logs in a parseable format (e.g., JSON or standardized text) containing timestamps, levels, and request context.

#### Scenario: Request logging
- **WHEN** an HTTP request is processed
- **THEN** the system logs the method, path, status, and duration in a consistent structure

### Requirement: Log Severity Levels
The system SHALL support different log levels (INFO, WARN, ERROR, DEBUG) to filter output based on environment.

#### Scenario: Error logging
- **WHEN** an internal server error occurs
- **THEN** it is logged with the ERROR level and includes the stack trace or detailed message
