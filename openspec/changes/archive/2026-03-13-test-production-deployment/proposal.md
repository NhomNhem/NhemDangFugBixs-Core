## Why

The backend has been refactored to Clean Architecture and deployed to Fly.io. We need to verify that the production environment is fully functional by performing a live integration test using a real PlayFab session token.

## What Changes

- **Production Testing**: Execute the integration test suite against the live `https://gamefeel-backend.fly.dev` endpoint.
- **Identity Verification**: Confirm that the backend correctly validates the provided PlayFab session ticket and returns a valid HW JWT.
- **End-to-End Validation**: Verify Save/Load and Leaderboard functionality in the real production environment (Supabase + Redis).

## Capabilities

### New Capabilities
- `production-verification`: Suite of tests specifically for the live environment.

### Modified Capabilities
- None.

## Impact

- **Reliability**: Confirms the system is ready for the Unity client.
- **Data Integrity**: Ensures production database and cache are correctly integrated with the new architecture.
