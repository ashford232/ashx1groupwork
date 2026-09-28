$ErrorActionPreference = "Stop"

$ROOT_DIR = Split-Path -Parent $PSScriptRoot
Set-Location $ROOT_DIR

# Load .env
if (Test-Path ".env") {
    Get-Content ".env" |
        Where-Object { $_ -and $_ -notmatch "^#" } |
        ForEach-Object {
            $key, $value = $_ -split "=", 2
            if ($key -and $value) {
                [Environment]::SetEnvironmentVariable($key.Trim(), $value.Trim())
            }
        }
}

if (-not $env:DATABASE_URL) {
    Write-Error "Error: DATABASE_URL is not set."
    exit 1
}

$MIGRATIONS_PATH = "migrations"
$ACTION = $args[0]
$STEPS = if ($args.Count -gt 1) { $args[1] } else { "1" }
$NAME = $args[1]

switch ($ACTION) {

    "up" {
        migrate -path $MIGRATIONS_PATH -database $env:DATABASE_URL up
    }

    "down" {
        migrate -path $MIGRATIONS_PATH -database $env:DATABASE_URL down $STEPS
    }

    "force" {
        migrate -path $MIGRATIONS_PATH -database $env:DATABASE_URL force $STEPS
    }

    "version" {
        migrate -path $MIGRATIONS_PATH -database $env:DATABASE_URL version
    }

    "create" {
        if (-not $NAME) {
            Write-Error "Error: Provide a migration name."
            exit 1
        }

        migrate create -ext sql -dir $MIGRATIONS_PATH -seq $NAME
    }

    default {
        Write-Host "Usage:"
        Write-Host "  .\scripts\migrate.ps1 up"
        Write-Host "  .\scripts\migrate.ps1 down [steps]"
        Write-Host "  .\scripts\migrate.ps1 force <version>"
        Write-Host "  .\scripts\migrate.ps1 version"
        Write-Host "  .\scripts\migrate.ps1 create <name>"
        exit 1
    }
}
