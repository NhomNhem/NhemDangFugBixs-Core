## Context

The project is a Go-based backend using the Fiber framework. It currently has a working health check and project structure but lacks automated build scripts and a comprehensive testing suite.

## Goals / Non-Goals

**Goals:**
- Automate the build process for multiple platforms.
- Establish a unit testing pattern using `testify`.
- Implement initial integration tests for core endpoints.
- Provide scripts for easy execution of build and test tasks.

**Non-Goals:**
- Setting up a full CI/CD pipeline (e.g., GitHub Actions) - this design focuses on the local and build-level tooling.
- Performance/Load testing.
- UI/Unity integration testing.

## Decisions

- **Testing Library**: Use `github.com/stretchr/testify` for assertions and mocking.
  - *Rationale*: It's the industry standard for Go, providing cleaner syntax than the standard library for common assertions.
- **Mocking Strategy**: Use `mockery` for generating mocks from interfaces.
  - *Rationale*: Manual mocking is error-prone and tedious. Automated generation ensures mocks stay in sync with interfaces.
- **Build Tooling**: Use a `Makefile` or Go scripts for build automation.
  - *Rationale*: Provides a consistent interface for developers and build systems.
- **Integration Testing**: Use `httptest` package from standard library to spin up a test server instance for API tests.
  - *Rationale*: Allows testing the full HTTP stack (middleware, routing, handlers) without requiring a separate running process.

## Risks / Trade-offs

- **Risk**: Mocking too much can lead to tests passing while the real system fails.
  - *Mitigation*: Ensure integration tests cover key paths with real (or semi-real) dependencies where possible.
- **Risk**: Environment differences between dev and build machines.
  - *Mitigation*: Use Docker for builds to ensure a consistent environment.
