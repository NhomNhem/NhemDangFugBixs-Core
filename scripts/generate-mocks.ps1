# Generate mocks using mockery
# Requirements: go install github.com/vektra/mockery/v2@latest

Write-Host "Updating mocks..." -ForegroundColor Cyan

# Clean existing mocks
if (Test-Path "internal/mocks") {
    Remove-Item -Recurse -Force "internal/mocks"
}

# Generate repository mocks
Write-Host "  Mocking repositories..."
go run github.com/vektra/mockery/v2 --dir internal/domain/repository --all --output internal/mocks/repository --outpkg repo_mock

# Generate usecase mocks
Write-Host "  Mocking usecases..."
go run github.com/vektra/mockery/v2 --dir internal/domain/usecase --all --output internal/mocks/usecase --outpkg usecase_mock

Write-Host "Mocks generated successfully!" -ForegroundColor Green
