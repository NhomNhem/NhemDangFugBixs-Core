# verify_build.ps1 - sanity check for backend

Write-Host "--- Starting Build Verification ---" -ForegroundColor Cyan

# 1. Compile Check
Write-Host "`n[1/3] Checking compilation..." -ForegroundColor Yellow
go build ./cmd/server
if ($LASTEXITCODE -ne 0) {
    Write-Host "Compilation failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}
Write-Host "✅ Compilation successful" -ForegroundColor Green

# 2. Unit & Integration Tests
Write-Host "`n[2/3] Running all tests..." -ForegroundColor Yellow
go test ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "Tests failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}
Write-Host "✅ All tests passed" -ForegroundColor Green

# 3. Swagger Generation Check
Write-Host "`n[3/3] Checking Swagger documentation..." -ForegroundColor Yellow
if (Get-Command swag -ErrorAction SilentlyContinue) {
    swag init -g cmd/server/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Swagger initialization failed!" -ForegroundColor Red
        exit $LASTEXITCODE
    }
    Write-Host "✅ Swagger documentation regenerated" -ForegroundColor Green
} else {
    Write-Host "⚠️  swag command not found, skipping Swagger check" -ForegroundColor Gray
}

Write-Host "`n🎉 Build verification complete!" -ForegroundColor Green
