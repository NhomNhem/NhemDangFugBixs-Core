## Context

The project has transitioned many features to Clean Architecture, but the documentation in `openspec/specs/` has not consistently kept pace. While some features have directories, many lack a formal `spec.md` or use incomplete placeholders. This design outlines the strategy for backfilling these specifications to ensure the entire system's behavior is formally documented and testable.

## Goals / Non-Goals

**Goals:**
- Provide a single source of truth for the behavior of all core systems.
- Standardize the specification format across the entire project.
- Ensure all requirements have at least one testable scenario.

**Non-Goals:**
- Implementation of new features (this change is for documentation of existing/planned behavior only).
- Refactoring of existing code (this will be handled by other changes like `full-architectural-refactor`).

## Decisions

### 1. Retroactive Specification Drafting
**Decision:** We will draft specifications based on the existing implementation and the architectural standards established during the Hollow Wilds phase 1.
**Rationale:** Since the code already exists for most of these systems, the specifications must accurately reflect the current "source of truth" while also incorporating the desired Clean Architecture standards.

### 2. Standardized Scenario Format
**Decision:** Strictly enforce the `#### Scenario: <name>` with `WHEN/THEN` format.
**Rationale:** Consistency in the specification layer allows for better automated processing and ensures that every requirement is articulated in a way that is directly translatable to a test case.

### 3. Capability-Centric Organization
**Decision:** Maintain the existing `openspec/specs/<capability>/` directory structure.
**Rationale:** This structure is already integrated with the OpenSpec CLI and provides a clear map of the system's modular capabilities.

## Risks / Trade-offs

- **[Risk] Specification-Code Mismatch** → **[Mitigation]** Use the `codebase_investigator` to verify that requirements match the current implementation in `internal/`.
- **[Trade-off] Effort vs. Code Velocity** → Writing detailed specs takes time away from coding, but it significantly reduces long-term technical debt and prevents "requirement drift."
