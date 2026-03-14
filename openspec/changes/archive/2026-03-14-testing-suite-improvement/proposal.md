## Why

The backend currently lacks a comprehensive automated testing suite, which increases the risk of regressions as new features are implemented. This change aims to establish formal testing patterns and initial coverage for critical paths.

## What Changes

- Implement a formal API integration testing framework.
- Create unit tests for core business logic in the `usecase` layer.
- Set up a repeatable pattern for mocking dependencies (database, Redis, external APIs).
- Add a build verification step to ensure the backend is always deployable.

## Capabilities

### New Capabilities
- `api-integration-testing`: Automated verification of HTTP endpoints using real request/response cycles.
- `unit-testing`: Isolation and verification of specific business rules without external dependencies.
- `build-verification`: Process to ensure the code compiles and passes basic sanity checks.

### Modified Capabilities
<!-- No requirement changes to existing capabilities yet -->

## Impact

- **Affected code**: All `internal/` packages will see the addition of `_test.go` files.
- **APIs**: Existing endpoints will be exercised by the integration suite.
- **Dependencies**: Adding `github.com/stretchr/testify` for assertions and mocking.
- **Systems**: Developer workflow will now include running the test suite locally.
