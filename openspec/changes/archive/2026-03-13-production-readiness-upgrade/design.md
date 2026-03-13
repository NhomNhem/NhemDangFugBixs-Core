## Context

The backend is currently in an "MVP" state regarding infrastructure. Migrations are manual, logging is basic, and the server process terminates abruptly. To reach production standards, we need to automate schema management and improve resource handling.

## Goals / Non-Goals

**Goals:**
- Automate database schema versioning.
- Implement parseable (JSON) logging for production.
- Ensure all database and cache connections are closed safely on exit.
- Restrict API access to authorized origins only.

**Non-Goals:**
- Zero-downtime migrations (we will accept brief downtime for schema changes in Phase 1).
- Distributed tracing (slated for a future phase).

## Decisions

### 1. Migration Management: Custom File-Based Runner
**Decision:** Implement a simple migration runner that scans the `migrations/` directory and tracks progress in a `schema_migrations` table.
**Rationale:** External libraries often add significant weight or complex dependencies. A simple, custom-built runner provides enough power for our current needs while remaining extremely easy to debug.
**Senior Detail:** Use a PostgreSQL Advisory Lock (e.g., `pg_advisory_lock(1234)`) to ensure that only one instance of the app executes migrations during a rolling deployment.

### 2. Structured Logging: `log/slog`
**Decision:** Transition all application logs to use the standard library's `log/slog` package.
**Rationale:** `slog` is the modern, performant, and official Go way to handle structured logging. It allows us to output JSON in production (for ingestion by Fly.io/CloudWatch) and readable text in development.
**Senior Detail:** Implement a `RequestID` middleware that injects a unique ID into the context and the logger, allowing for correlation of logs across a single player request.

### 3. Graceful Shutdown: Signal Context
**Decision:** Use `os/signal` to listen for SIGINT/SIGTERM and a 10-second timeout context to drain connections.
**Rationale:** This prevents data loss or "orphaned" connections in Supabase/Redis when the Fly.io machine restarts or scales.
**Senior Detail:** Orchestrate the shutdown: first `app.Shutdown()` to stop accepting traffic, then a short drain period, and finally closing the database and Redis pools.

### 4. CORS Security: Environment-Based Allowlist
**Decision:** Implement a CORS middleware that pulls allowed origins from an `ALLOWED_ORIGINS` environment variable.
**Rationale:** Prevents unauthorized websites from making requests to the API, while maintaining flexibility across development, staging, and production.

## Risks / Trade-offs

- **[Risk] Migration Deadlock** → **[Mitigation]** Use a row-level lock on the `schema_migrations` table during execution to prevent concurrent startups from interfering.
- **[Trade-off] Logging Overhead** → JSON logging is slightly more verbose than standard text. **[Mitigation]** This is offset by the ease of querying and monitoring in production.
