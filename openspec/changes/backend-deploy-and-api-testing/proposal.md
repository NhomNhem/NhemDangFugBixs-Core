## Why

The team needs a repeatable way to validate Docker-based deployment and production API behavior before and after release. Current checks are partly manual and inconsistent, which increases the risk of shipping broken backend deployments.

## What Changes

- Define a standardized deployment verification flow for Docker build, deploy, and post-deploy checks.
- Define a production smoke test flow for critical API endpoints against the real backend URL.
- Add clear pass/fail criteria for health, auth, and core gameplay endpoints during release verification.
- Align CI/CD and manual deployment validation so both paths enforce the same checks.

## Capabilities

### New Capabilities
- `deployment-and-api-verification`: End-to-end verification requirements for Docker deploy readiness, deployment execution, and real-backend API smoke testing.

### Modified Capabilities
- None.

## Impact

- Affected code: deployment scripts/workflows, test scripts, and verification docs.
- APIs: health, auth, player, and leaderboard endpoints used in smoke tests.
- Systems: Docker image build pipeline, Fly.io deployment path, production backend environment.
