## Why

The project currently lacks a comprehensive automated testing suite. While we have integration scripts, we need proper unit and integration testing integrated into the development workflow to ensure regressions are caught early and maintain high code quality.

## What Changes

- **Test Infrastructure**: Set up a standardized testing framework using Go's built-in testing package along with `testify` for assertions and `mockery` for mock generation.
- **Unit Tests**: Implement unit tests for core usecases (`Auth`, `Player`, `Leaderboard`) using mocks for repositories.
- **Integration Tests**: Refactor the existing test scripts into Go test files (`*_test.go`) that can be run with `go test ./...`.
- **CI/CD Integration**: Add a testing step to the GitHub Actions workflow to run all tests before deployment.

## Capabilities

### New Capabilities
- `unit-testing`: Standardized patterns for testing business logic in isolation.
- `api-integration-testing`: Automated verification of HTTP endpoints against a test database.
- `mock-generation`: Automated creation of repository and usecase mocks.

### Modified Capabilities
- `deployment`: Added requirement for successful test execution before production pushes.

## Impact

- **Development Speed**: Faster feedback loops for developers.
- **Reliability**: significantly reduced risk of breaking existing features during refactoring.
- **Documentation**: Tests serve as live documentation for system behavior.
