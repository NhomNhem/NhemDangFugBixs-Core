# generate-mocks.ps1 - Automated mock generation

Write-Host "Updating mocks..." -ForegroundColor Cyan

# Ensure internal/mocks directory exists
if (-not (Test-Path internal/mocks)) {
    New-Item -ItemType Directory -Path internal/mocks
}

# Run mockery
mockery

Write-Host "Mocks generated successfully!" -ForegroundColor Green
