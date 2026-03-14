## 1. Pre-Deployment Checks

- [x] 1.1 Run `verify_build.ps1` locally to ensure all tests pass
- [x] 1.2 Verify `fly.toml` configuration is up to date

## 2. Deployment

- [x] 2.1 Trigger deployment to Fly.io (via GitHub Action or manual `fly deploy`)
- [x] 2.2 Monitor deployment logs for any startup errors

## 3. Production Verification

- [x] 3.1 Verify `GET /health` returns status `ok` and connected database
- [x] 3.2 Verify `GET /swagger/index.html` is accessible
- [x] 3.3 Run smoke tests for critical endpoints (Auth, Leaderboard) against production URL

## 4. Finalization

- [x] 4.1 Update `DEPLOYMENT_CHECKLIST.md` with any new verification steps
