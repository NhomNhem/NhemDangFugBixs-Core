# build.ps1 - Cross-platform build script for Hollow Wilds Backend

$VERSION = "1.0.0"
$COMMIT = git rev-parse --short HEAD
$LDFLAGS = "-X main.Version=$VERSION -X main.Commit=$COMMIT"

# Create bin directory
if (-not (Test-Path bin)) {
    New-Item -ItemType Directory -Path bin
}

# Windows Build
Write-Host "Building for Windows..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o bin/server.exe cmd/server/main.go

# Linux Build
Write-Host "Building for Linux..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o bin/server cmd/server/main.go

# Darwin Build (MacOS)
Write-Host "Building for MacOS..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o bin/server-darwin cmd/server/main.go

Write-Host "Builds complete. Binaries located in bin/" -ForegroundColor Green
