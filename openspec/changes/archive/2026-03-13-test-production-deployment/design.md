## Context

The backend has been refactored and deployed to production. We have a real PlayFab session token provided by the user to perform a final validation of the live environment.

## Goals / Non-Goals

**Goals:**
- Verify the production endpoint `https://gamefeel-backend.fly.dev` is fully operational.
- Validate real PlayFab ticket authentication.
- Confirm production database persistence.

**Non-Goals:**
- Local testing (this is specifically for production).
- Permanent storage of the session token.

## Decisions

### 1. Test Execution Method
**Decision:** Modify `scripts/test_hollow_wilds.go` to accept a base URL and token from the environment.
**Rationale:** Reuse existing test logic while allowing flexible configuration for production verification.

### 2. Verified Data
**Decision:** Use the PlayFab ID `8B7957C000B402F0` and its associated token for the test.
**Rationale:** These were provided by the user from a successful Unity client login.

## Risks / Trade-offs

- **[Risk] Production Data Pollution** → **[Mitigation]** The test will use a specific test player ID and its data can be manually cleaned up in Supabase if necessary.
- **[Risk] Token Expiration** → **[Mitigation]** Perform the test promptly before the ticket expires (PlayFab tickets usually last 24 hours).
