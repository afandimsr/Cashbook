param (
    [Parameter(Mandatory=$true)]
    [string]$name
)

$timestamp = Get-Date -Format "yyyyMMddHHmmss"
$upFile = "migrations/${timestamp}_${name}.up.sql"
$downFile = "migrations/${timestamp}_${name}.down.sql"

New-Item -ItemType Directory -Force -Path "migrations" | Out-Null
New-Item -ItemType File -Force -Path $upFile | Out-Null
New-Item -ItemType File -Force -Path $downFile | Out-Null

Write-Host "Created $upFile"
Write-Host "Created $downFile"
