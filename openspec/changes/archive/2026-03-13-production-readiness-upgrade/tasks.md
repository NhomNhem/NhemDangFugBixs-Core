## 1. Automated Migration Manager

- [x] 1.1 Create `internal/database/migrator.go` to handle schema versioning
- [x] 1.2 Implement `RunMigrations` logic with PostgreSQL Advisory Locks for safety
- [x] 1.3 Add `schema_migrations` tracking table to ensure idempotency
- [x] 1.4 Integrate the migrator into the server startup sequence in `main.go`

## 2. Structured Logging & Context

- [x] 2.1 Initialize a global `slog` handler in `main.go` (JSON for production, Text for dev)
- [x] 2.2 Create `RequestID` middleware to correlate logs across a single request
- [x] 2.3 Refactor `internal/database` to use structured logs for connection events
- [x] 2.4 Refactor `internal/delivery` handlers to include request context in logs

## 3. Graceful Shutdown

- [x] 3.1 Implement a signal listener in `main.go` for SIGINT and SIGTERM
- [x] 3.2 Create an orchestrated shutdown flow: Fiber -> Drain -> DB -> Redis
- [x] 3.3 Ensure a timeout context is used to prevent "hanging" shutdowns

## 4. Security Hardening

- [x] 4.1 Update CORS configuration to restrict origins via `ALLOWED_ORIGINS` env var
- [x] 4.2 Add a startup check to validate that `JWT_SECRET` and `DATABASE_URL` meet length/complexity requirements
- [x] 4.3 Ensure `PLAYFAB_TITLE_ID` is never empty in production environments

