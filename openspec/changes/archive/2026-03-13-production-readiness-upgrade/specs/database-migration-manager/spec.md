## ADDED Requirements

### Requirement: Automated Migration Execution
The system SHALL automatically apply all pending SQL migrations from the `migrations/` directory upon server startup.

#### Scenario: Successful migration on startup
- **WHEN** the server starts and finds new `.sql` files in the migrations folder
- **THEN** it executes them in order and records the successful application in a tracking table

### Requirement: Migration Version Control
The system SHALL maintain a `schema_migrations` table to track which files have already been executed.

#### Scenario: No pending migrations
- **WHEN** the server starts and all migration files match the tracking table
- **THEN** it proceeds to start the API without executing any SQL
