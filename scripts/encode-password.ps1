Add-Type -AssemblyName System.Web

Write-Host "=== PASSWORD URL ENCODER ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Paste your FULL Supabase connection string:" -ForegroundColor Yellow
$connStr = Read-Host

if ($connStr -match 'postgres://postgres:(.+?)@(.+)' -or $connStr -match 'postgresql://postgres:(.+?)@(.+)') {
    $pwd = $matches[1]
    $rest = $matches[2]
    
    $encoded = [System.Web.HttpUtility]::UrlEncode($pwd)
    $newConn = "postgresql://postgres:$encoded@$rest"
    
    Write-Host ""
    Write-Host "Original password: $pwd" -ForegroundColor White
    Write-Host "Encoded password:  $encoded" -ForegroundColor Green
    Write-Host ""
    Write-Host "NEW CONNECTION STRING:" -ForegroundColor Cyan
    Write-Host $newConn -ForegroundColor Yellow
    Write-Host ""
    
    Set-Clipboard -Value $newConn
    Write-Host "Copied to clipboard! Paste to .env file" -ForegroundColor Green
} else {
    Write-Host "ERROR: Invalid format" -ForegroundColor Red
}

Read-Host "Press Enter to exit"
