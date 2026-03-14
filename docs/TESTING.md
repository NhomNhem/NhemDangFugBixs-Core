# Testing Guide

This document outlines the testing strategy and patterns for the Hollow Wilds Backend.

## Testing Strategy

We employ a two-tiered testing strategy:

1.  **Unit Tests**: Focused on individual business rules in the `usecase` layer. These tests use mocks for all external dependencies (repositories, identity providers, etc.).
2.  **Integration Tests**: Focused on the API layer (`internal/api`). These tests exercise the full HTTP request/response cycle, including routing and middleware, using mocked usecases.

## Unit Testing

Unit tests are located alongside the code they test (e.g., `internal/usecase/auth/auth_usecase_test.go`).

### Mocking

We use [Mockery](https://github.com/vektra/mockery) to generate mocks for our interfaces. 

#### Generating Mocks
To update all mocks after an interface change, run:
```powershell
./scripts/generate-mocks.ps1
```

#### Using Mocks in Tests
We use `testify/mock` for behavior-based verification.
```go
repo := new(repository_mock.MockPlayerRepository)
repo.On("GetByID", ctx, playerID).Return(expectedPlayer, nil).Once()
```

## Integration Testing

Integration tests verify that our handlers, routes, and middleware work together correctly. We use Fiber's built-in `Test` method which uses `httptest` under the hood without opening real network ports.

### Test Helpers
Use `SetupTestApp()` in `internal/api/test_helpers.go` to get a pre-configured Fiber app with all usecases mocked.

```go
app, mocks := SetupTestApp()
mocks.Auth.On("Login", ...).Return(...)

req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
resp, err := app.Test(req)
```

## Build Verification

Before every commit, it is recommended to run the build verification script to ensure compilation, tests, and documentation are all in sync:

```powershell
./scripts/verify_build.ps1
```
