## Context

The backend has successfully moved its core survival features (Hollow Wilds) to Clean Architecture. However, several legacy systems—Admin, Talents, and Levels—still reside in the old `internal/api` and `internal/services` structure. These older services access the database directly via global variables and lack proper testability.

## Goals / Non-Goals

**Goals:**
- Move all remaining business logic from legacy services to the Usecase layer.
- Encapsulate all database access for these systems within Repositories.
- Unified Dependency Injection across the entire project.
- Safely delete the `internal/api` and `internal/services` directories.

**Non-Goals:**
- Database schema changes (this is a structural refactor only).
- API endpoint changes (routes and response bodies must remain identical).

## Decisions

### 1. Interface-First Migration
**Decision:** We will define `AdminRepository`, `TalentRepository`, and `LevelRepository` interfaces in the Domain layer before moving any implementation code.
**Rationale:** This ensures we maintain the Dependency Rule throughout the refactor and allows us to generate mocks immediately.

### 2. Module-by-Module Refactor
**Decision:** Refactor Talents, then Levels, and finally Admin.
**Rationale:** The Admin system is the largest and most complex; starting with smaller modules allows us to verify the transition process more easily.

### 3. Unified Request Correlation
**Decision:** Ensure the new usecases accept `context.Context` and leverage the `requestId` from the middleware for logging.
**Rationale:** Maintains consistency with the high standards set in the previous production-readiness upgrade.

## Risks / Trade-offs

- **[Risk] Broken Imports** → **[Mitigation]** Use the `generalist` agent to perform project-wide import updates and verify with `go build`.
- **[Trade-off] Increased File Count** → Clean Architecture inherently involves more files due to interface definitions. This is a trade-off for significantly better long-term maintainability.
