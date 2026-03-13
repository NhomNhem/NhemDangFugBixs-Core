# check-quality.ps1 - Lint and Test runner

Write-Host "--- Running Quality Checks ---" -ForegroundColor Cyan

# 1. Run Linter
Write-Host "`n[1/2] Running Linter..." -ForegroundColor Yellow
if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
    golangci-lint run
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Linter failed!" -ForegroundColor Red
        exit $LASTEXITCODE
    }
} else {
    Write-Host "golangci-lint not found, skipping linting step." -ForegroundColor Gray
}

# 2. Run Tests
Write-Host "`n[2/2] Running Tests..." -ForegroundColor Yellow
go test -v ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "Tests failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "`n✅ All quality checks passed!" -ForegroundColor Green
