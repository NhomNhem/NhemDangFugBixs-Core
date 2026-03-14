## ADDED Requirements

### Requirement: Health Endpoint Accessibility
The production environment SHALL provide a public `/health` endpoint that returns a `200 OK` status when the application and its dependencies are healthy.

#### Scenario: Successful production health check
- **WHEN** a client performs a `GET https://gamefeel-backend.fly.dev/health`
- **THEN** the system returns `200 OK` with status `ok` and version information

### Requirement: API Documentation Availability
The production environment SHALL serve the Swagger UI at `/swagger/index.html`.

#### Scenario: Accessing live documentation
- **WHEN** a developer visits `https://gamefeel-backend.fly.dev/swagger/index.html`
- **THEN** the Swagger UI is rendered correctly with the latest API definitions

### Requirement: Critical Path Verification
The production environment SHALL successfully process core player requests (Auth, Save, Leaderboard).

#### Scenario: Verifying login in production
- **WHEN** a client submits a valid PlayFab ticket to `/api/v1/auth/hw/login`
- **THEN** the system returns a valid JWT and player information
