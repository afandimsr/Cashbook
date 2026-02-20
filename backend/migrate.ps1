param (
    [Parameter(Mandatory = $true)]
    [ValidateSet("up", "down", "force", "swag")]
    [string]$Action,

    [Parameter(Mandatory = $false)]
    [string]$Version,

    [Parameter(Mandatory = $false)]
    [switch]$ConfirmInput
)

# Load .env file
if (Test-Path .env) {
    Get-Content .env | Where-Object { $_ -match '=' } | ForEach-Object {
        $key, $value = $_.Split('=', 2)
        [Environment]::SetEnvironmentVariable($key.Trim(), $value.Trim(), "Process")
    }
}
else {
    Write-Error ".env file not found"
    exit 1
}

# Construct Database URL
if ($env:DATABASE_URL) {
    $MIGRATE_URL = $env:DATABASE_URL
}
else {
    $DB_DRIVER = $env:DB_DRIVER
    if ($null -eq $DB_DRIVER) { $DB_DRIVER = "postgres" }
    
    $DB_USER = $env:DB_USER
    $DB_PASS = $env:DB_PASS
    $DB_HOST = $env:DB_HOST
    $DB_PORT = $env:DB_PORT
    $DB_NAME = $env:DB_NAME

    if ($DB_DRIVER -eq "mysql") {
        $MIGRATE_URL = "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"
    }
    elseif ($DB_DRIVER -eq "postgres") {
        $MIGRATE_URL = "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
    }
    else {
        Write-Error "Unsupported DB_DRIVER: $DB_DRIVER. Please use 'mysql' or 'postgres', or set DATABASE_URL in .env"
        exit 1
    }
}

$APP_ENV = $env:APP_ENV

# Safeguard against accidental production rollback
if ($Action -eq "down" -and $APP_ENV -eq "production" -and -not $ConfirmInput) {
    Write-Host "ERROR: DANGER! You are attempting to rollback in PRODUCTION." -ForegroundColor Red
    Write-Host "This can cause DATA LOSS. To proceed, use: .\migrate.ps1 -Action down -ConfirmInput" -ForegroundColor Yellow
    exit 1
}

# Execute migrate command
$Cmd = "migrate -database ""$MIGRATE_URL"" -path migrations"

if ($Action -eq "swag") {
    swag init -g cmd/api/main.go -o docs
}
elseif ($Action -eq "force") {
    if (-not $Version) {
        Write-Error "Version parameter is required for force action"
        exit 1
    }
    Invoke-Expression "$Cmd force $Version"
}
else {
    Invoke-Expression "$Cmd $Action"
}
