## 1. Deployment Validation Setup

- [x] 1.1 Review and align Docker build/deploy scripts or workflow steps for release verification
- [x] 1.2 Define required runtime configuration checks (env/secrets/ports) before deployment

## 2. Deployment Verification Flow

- [x] 2.1 Add deployment completion checks for CI/CD and manual deploy paths
- [x] 2.2 Add startup log/health verification step to fail release on startup errors

## 3. Production API Smoke Tests

- [x] 3.1 Add smoke test cases for `/health` and `/swagger/index.html` against production URL
- [x] 3.2 Add smoke test cases for critical endpoints (auth, player, leaderboard) with expected responses

## 4. Release Gating and Documentation

- [x] 4.1 Wire verification checks into release gate criteria with pass/fail outcomes
- [x] 4.2 Update deployment/testing documentation with the standardized verification runbook
