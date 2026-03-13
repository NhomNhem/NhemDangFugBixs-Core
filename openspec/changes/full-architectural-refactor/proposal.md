## Why

The backend currently operates in a "hybrid" state where some features follow the new Clean Architecture (Domain -> Usecase -> Infrastructure) while others (Admin, Talents, Levels) remain in the legacy tightly-coupled structure. This creates technical debt, makes testing inconsistent, and confuses developers. We need a 100% unified architecture.

## What Changes

- **Legacy Removal**: Fully eliminate the `internal/api` and `internal/services` directories.
- **Admin Refactor**: Move Admin logic to `internal/usecase/admin` and `internal/infrastructure/persistence`.
- **Talent Refactor**: Move Talent logic to `internal/usecase/talent` and `internal/infrastructure/persistence`.
- **Level Refactor**: Move Level logic to `internal/usecase/level` and `internal/infrastructure/persistence`.
- **Unified Dependency Injection**: Ensure all remaining handlers in `internal/delivery/http` use constructor injection.

## Capabilities

### New Capabilities
- `admin-management`: Secure administrative operations refactored into Clean Architecture.
- `talent-system`: Character talent progression and upgrades.
- `level-progression`: Level completion tracking and statistics.

### Modified Capabilities
- None (The requirements remain the same; only the internal structure is changing to match the project's new standards).

## Impact

- **Codebase Cleanliness**: Zero legacy "flat" services remaining.
- **Testability**: 100% of business logic becomes unit-testable with mocks.
- **Maintainability**: A single, predictable pattern for all developers to follow.
