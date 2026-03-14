## Context

The backend is deployed with Docker artifacts and hosted on Fly.io. Verification currently mixes ad-hoc manual checks and partially automated scripts, which makes releases harder to trust. This change defines one consistent verification flow covering Docker image validation, deployment execution, and production API smoke tests against the real backend.

## Goals / Non-Goals

**Goals:**
- Standardize pre-deploy and post-deploy verification for Docker-based releases.
- Define mandatory production smoke tests for critical API paths.
- Ensure CI/CD-triggered and manual deploys use equivalent validation criteria.
- Provide clear go/no-go gates and rollback signals.

**Non-Goals:**
- Performance, load, or stress testing.
- Redesign of backend domain logic.
- Migration to a different hosting provider or container runtime.

## Decisions

- **Single verification contract for deploy paths**: Use one requirement set for both manual and CI/CD deploy flows so release quality does not depend on deployment method.
  - *Alternative considered*: Separate requirements per deploy path.
  - *Why not*: Creates drift and duplicated maintenance.
- **Production smoke tests focus on critical endpoints**: Require health, auth, player, and leaderboard verification against real production URL.
  - *Alternative considered*: Full regression suite in production.
  - *Why not*: Too slow and risky for release gates.
- **Fail-fast deployment validation**: If Docker build/deploy verification or critical endpoint checks fail, release is treated as failed and rollback procedures are initiated.
  - *Alternative considered*: Allow non-critical failures and continue.
  - *Why not*: Increases probability of user-facing incidents.

## Risks / Trade-offs

- **Risk: Overly strict checks block releases** → Mitigation: scope mandatory checks to critical paths only.
- **Risk: Production-only issues still missed** → Mitigation: require real-backend smoke tests immediately after deploy and monitor logs.
- **Risk: Verification maintenance overhead** → Mitigation: keep endpoint set small and tie checks to stable API contracts.

## Migration Plan

1. Document and adopt the new verification checklist in deployment workflow.
2. Update CI/CD steps to run the same Docker and API smoke checks.
3. Run the standardized flow for upcoming releases and tune thresholds if needed.
4. If deployment fails verification, rollback to previous healthy release and log the failed step for remediation.

## Open Questions

- Should smoke tests run only post-deploy or also against a staged environment pre-deploy?
- Which endpoint-level assertions should be considered mandatory vs informational as the API evolves?
