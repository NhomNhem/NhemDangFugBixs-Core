## ADDED Requirements

### Requirement: Multi-Environment CI/CD
The system SHALL support automated deployment to distinct environments (Development, Staging, Production) based on branch merges and tags.

#### Scenario: Merging to main
- **WHEN** a pull request is merged into the `main` branch
- **THEN** the CI/CD pipeline runs tests and deploys the backend to the `staging` environment

### Requirement: Fly.io Integration
The system SHALL use Fly.io for container orchestration and global traffic management.

#### Scenario: Production release
- **WHEN** a new git tag matching `v*.*.*` is pushed
- **THEN** the pipeline builds a Docker image and performs a rolling update on Fly.io
