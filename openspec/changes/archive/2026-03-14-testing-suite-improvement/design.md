## Context

The backend is built with Go using the Fiber framework. It follows a Clean Architecture pattern with `delivery/http` (handlers), `usecase` (business logic), and `infrastructure/persistence` (repositories). Currently, testing is manual or ad-hoc.

## Goals / Non-Goals

**Goals:**
- Implement a suite of unit tests for the `usecase` layer.
- Implement an end-to-end integration test suite for the `api` layer.
- Standardize dependency injection to support mocking.
- Automate the build verification process.

**Non-Goals:**
- 100% code coverage (focus is on critical paths).
- Performance/Load testing.
- Frontend/Unity client side testing.

## Decisions

- **Assertion Library**: Use `github.com/stretchr/testify/assert`.
  - *Rationale*: Provides a much cleaner and more readable syntax than standard library `if err != nil` checks.
- **Mocking Strategy**: Use `github.com/stretchr/testify/mock`.
  - *Rationale*: Allows for behavior-based verification and is well-integrated with the assertion library.
- **Integration Test Execution**: Use `httptest` and Fiber's `app.Test`.
  - *Rationale*: Fiber provides a built-in `Test` method that allows exercising the full stack (routing, middleware, handlers) without opening a real network port.
- **Test Database**: Use mock repositories for unit tests; use a dedicated local test database (or transaction rollback) for integration tests if needed.
  - *Rationale*: Keeps unit tests fast and deterministic.

## Risks / Trade-offs

- **Risk**: Brittle tests that break on every refactor.
  - *Mitigation*: Focus unit tests on the `usecase` layer interfaces rather than internal implementation details.
- **Risk**: Inconsistent environments leading to "works on my machine" failures.
  - *Mitigation*: Provide a Docker-based test environment and include build verification in the CI process.
