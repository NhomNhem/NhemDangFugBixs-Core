## 1. Pre-Deployment Checks

- [ ] 1.1 Perform a full backup of the production Supabase database
- [x] 1.2 Review `fly.toml` to ensure correct internal port and resource allocation
- [x] 1.3 Align `analytics_events` schema and code with production context

## 2. Secrets & CI/CD Configuration

- [x] 2.1 Set production database secret: `fly secrets set SUPABASE_DATABASE_URL="..."` (Configured via GitHub Actions)
- [x] 2.2 Set production Redis secret: `fly secrets set REDIS_URL="..."` (Configured via GitHub Actions)
- [x] 2.3 Set JWT security secret: `fly secrets set JWT_SECRET="..."` (Configured via GitHub Actions)
- [x] 2.4 Set PlayFab Title ID: `fly secrets set PLAYFAB_TITLE_ID="..."` (Configured via GitHub Actions)
- [x] 2.5 Configure GitHub Secret: Add `FLY_API_TOKEN` to repository secrets

## 3. Production Migration

- [x] 3.1 Execute `004_hollow_wilds_phase1.sql` migration against the production database
- [x] 3.2 Verify table creation (`players`, `player_saves`, etc.) in the Supabase Dashboard

## 4. Deployment & Verification

- [ ] 4.1 Trigger automated deployment: Push latest changes to `main` branch
- [ ] 4.2 Monitor GitHub Actions workflow for successful completion
- [ ] 4.3 Verify live health check status: `https://gamefeel-backend.fly.dev/health`
- [ ] 4.4 Perform a production smoke test using `scripts/test_hollow_wilds.go` (pointing to the live URL)
