Write-Host "=== Password Encoder ===" -ForegroundColor Cyan
Write-Host "Paste your password and press Enter:"
$password = Read-Host
Add-Type -AssemblyName System.Web
$encoded = [System.Web.HttpUtility]::UrlEncode($password)
Write-Host ""
Write-Host "ENCODED PASSWORD:" -ForegroundColor Green
Write-Host $encoded
Write-Host ""
Write-Host "Copy this encoded password to your .env file" -ForegroundColor Yellow
