## 1. Domain & Infrastructure (Repositories)

- [x] 1.1 Define `AdminRepository`, `TalentRepository`, and `LevelRepository` interfaces in `internal/domain/repository`
- [x] 1.2 Implement `PostgresTalentRepository` in `internal/infrastructure/persistence`
- [x] 1.3 Implement `PostgresLevelRepository` in `internal/infrastructure/persistence`
- [x] 1.4 Implement `PostgresAdminRepository` in `internal/infrastructure/persistence`

## 2. Usecase Layer Refactor

- [x] 2.1 Implement `TalentUsecase` using the new repository
- [x] 2.2 Implement `LevelUsecase` using the new repository
- [x] 2.3 Implement `AdminUsecase` using the new repository
- [x] 2.4 Regenerate all mocks using `scripts/generate-mocks.ps1`

## 3. Delivery Layer (Handlers)

- [x] 3.1 Create `internal/delivery/http/talent_handler.go` with constructor injection
- [x] 3.2 Create `internal/delivery/http/level_handler.go` with constructor injection
- [x] 3.3 Create `internal/delivery/http/admin_handler.go` with constructor injection

## 4. Main Integration & Cleanup

- [x] 4.1 Initialize new repositories and usecases in `main.go`
- [x] 4.2 Update route registration to use the new handlers
- [x] 4.3 Delete `internal/api` and `internal/services` directories
- [x] 4.4 Verify build and run existing integration tests: `go run scripts/test_hollow_wilds.go`

## 5. Testing

- [x] 5.1 Add unit tests for `TalentUsecase`
- [x] 5.2 Add unit tests for `LevelUsecase`
- [x] 5.3 Add unit tests for `AdminUsecase`
