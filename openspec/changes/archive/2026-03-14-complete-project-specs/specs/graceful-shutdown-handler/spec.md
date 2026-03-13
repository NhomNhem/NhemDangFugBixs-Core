## ADDED Requirements

### Requirement: Signal Handling
The system SHALL intercept `SIGINT` (Ctrl+C) and `SIGTERM` signals to initiate a controlled shutdown sequence.

#### Scenario: Receiving termination signal
- **WHEN** the process receives a `SIGTERM` from the OS or orchestrator
- **THEN** it stops accepting new requests and waits for in-flight requests to complete

### Requirement: Resource Cleanup
The system SHALL ensure that all persistent connections (Database, Redis, etc.) are closed properly during shutdown.

#### Scenario: Closing database pool
- **WHEN** the server is in the shutdown phase
- **THEN** the `pgxpool.Pool` is closed after all active queries have finished
