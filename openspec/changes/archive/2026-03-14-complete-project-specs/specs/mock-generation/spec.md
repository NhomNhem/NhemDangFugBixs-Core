## ADDED Requirements

### Requirement: Automated Repository Mocking
The system SHALL use `mockery` to automatically generate and maintain mock implementations for all repository and usecase interfaces.

#### Scenario: Interface update
- **WHEN** a developer adds a new method to a repository interface in `internal/domain/repository`
- **THEN** running the `generate-mocks.ps1` script creates the corresponding mock in `internal/mocks/repository`

### Requirement: Unified Test Double Pattern
All unit tests SHALL use the generated mocks with standard `Expect` and `Return` behaviors.

#### Scenario: Testing a usecase
- **WHEN** a developer creates a new unit test for a usecase
- **THEN** they inject the generated mock implementation into the usecase constructor
