## 1. Domain & Infrastructure (Repositories)

- [ ] 1.1 Define `AdminRepository`, `TalentRepository`, and `LevelRepository` interfaces in `internal/domain/repository`
- [ ] 1.2 Implement `PostgresTalentRepository` in `internal/infrastructure/persistence`
- [ ] 1.3 Implement `PostgresLevelRepository` in `internal/infrastructure/persistence`
- [ ] 1.4 Implement `PostgresAdminRepository` in `internal/infrastructure/persistence`

## 2. Usecase Layer Refactor

- [ ] 2.1 Implement `TalentUsecase` using the new repository
- [ ] 2.2 Implement `LevelUsecase` using the new repository
- [ ] 2.3 Implement `AdminUsecase` using the new repository
- [ ] 2.4 Regenerate all mocks using `scripts/generate-mocks.ps1`

## 3. Delivery Layer (Handlers)

- [ ] 3.1 Create `internal/delivery/http/talent_handler.go` with constructor injection
- [ ] 3.2 Create `internal/delivery/http/level_handler.go` with constructor injection
- [ ] 3.3 Create `internal/delivery/http/admin_handler.go` with constructor injection

## 4. Main Integration & Cleanup

- [ ] 4.1 Initialize new repositories and usecases in `main.go`
- [ ] 4.2 Update route registration to use the new handlers
- [ ] 4.3 Delete `internal/api` and `internal/services` directories
- [ ] 4.4 Verify build and run existing integration tests: `go run scripts/test_hollow_wilds.go`

## 5. Testing

- [ ] 5.1 Add unit tests for `TalentUsecase`
- [ ] 5.2 Add unit tests for `LevelUsecase`
- [ ] 5.3 Add unit tests for `AdminUsecase`
