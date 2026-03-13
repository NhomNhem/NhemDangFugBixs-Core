## Context

The backend has recently been refactored into a Clean Architecture structure. This reorganization provides clear boundaries (Domain, Usecase, Infrastructure, Delivery) which are ideal for automated testing. Currently, we only have external verification scripts; we need to formalize this into a Go-native testing suite.

## Goals / Non-Goals

**Goals:**
- Implement unit tests for all Usecases using interface-based mocks.
- Implement integration tests for core HTTP endpoints.
- Automate mock generation to reduce developer boilerplate.
- Enforce testing in the CI/CD pipeline.

**Non-Goals:**
- 100% code coverage (we will prioritize critical business paths first).
- Testing legacy "GameFeel" generic handlers (focus is on Hollow Wilds).

## Decisions

### 1. Mocking Tool: Mockery
**Decision:** Use `github.com/vektra/mockery/v2`.
**Rationale:** It is the industry standard for Go. It generates type-safe mocks from interfaces automatically, making unit testing the Usecase layer much easier.

### 2. Assertion Library: Testify
**Decision:** Use `github.com/stretchr/testify`.
**Rationale:** Provides clear, readable assertions (`assert.Equal`, `assert.NoError`) which are superior to standard library `if err != nil` checks in test code.

### 3. Integration Testing: Fiber's Test App
**Decision:** Use Fiber's built-in `app.Test()` method.
**Rationale:** It allows running full HTTP requests against the router in-memory, providing high-fidelity integration testing without needing to manage network ports.

### 4. CI/CD Enforcement
**Decision:** Add a `go test ./...` step to `.github/workflows/deploy.yml`.
**Rationale:** Ensures that no code is deployed to Fly.io if the testing suite fails.

## Risks / Trade-offs

- **[Risk] Brittle Tests** → If mocks are too detailed, tests may break on minor logic changes. **[Mitigation]** Focus on verifying behaviors and state transitions rather than implementation details.
- **[Trade-off] Build Time** → Adding tests to CI will increase the time it takes to deploy. **[Mitigation]** Use Go's build cache and only run tests on relevant changes.
