## Why

We need to make the Hollow Wilds Phase 1 backend features available for the Unity client in a production environment. This transition ensures that authentication, save/load, and leaderboard functionalities are accessible beyond local development.

## What Changes

- **Fly.io Deployment**: Push the latest containerized backend to Fly.io infrastructure.
- **Secrets Management**: Configure production-grade secrets (DATABASE_URL, REDIS_URL, JWT_SECRET, PLAYFAB_TITLE_ID) on Fly.io and GitHub.
- **GitHub Actions Workflow**: Update and validate the `.github/workflows/deploy.yml` for automated production pushes.
- **Production Migrations**: Execute the `004_hollow_wilds_phase1.sql` migration against the production Supabase instance.
- **Endpoint Update**: Finalize the base URL for the production API.

## Capabilities

### New Capabilities
- `deployment`: Procedures for deploying and scaling the backend on Fly.io.
- `production-ready`: Final validation of connectivity, security, and schema integrity in the live environment.

### Modified Capabilities
- None.

## Impact

- **Infrastructure**: Fly.io application will be updated with the new Phase 1 codebase.
- **Database**: Production Supabase schema will be updated to include `players`, `player_saves`, `player_save_backups`, and `leaderboard_entries`.
- **Unity Client**: Will transition from local testing to the production endpoint `api.hollowwilds.com`.
