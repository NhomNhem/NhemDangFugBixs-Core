## Context

The backend is currently deployed on Fly.io. We have recently implemented production readiness upgrades (standardized errors, rate limiting, DB configs) and need to ensure these are correctly deployed and verified in the live environment.

## Goals / Non-Goals

**Goals:**
- Deploy latest binary to Fly.io.
- Verify environmental variables and secrets are correctly set.
- Run integration tests against the live URL.
- Confirm Swagger UI is accessible.

**Non-Goals:**
- Performance or Load testing (handled separately).
- Changes to infrastructure code (fly.toml).

## Decisions

- **Verification Tool**: Use `curl` and the project's own integration tests (adapted for base URL) to verify health.
- **Deployment Method**: Use `fly deploy` via the generalist subagent if possible, or manual push.

## Risks / Trade-offs

- **Risk**: Production downtime during deployment.
  - *Mitigation*: Fly.io handles rolling updates by default.
- **Risk**: Secret mismatch in production.
  - *Mitigation*: Run health checks immediately after deployment to verify DB/Redis connectivity.
