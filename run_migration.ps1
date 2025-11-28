# PowerShell script to run database migrations
# Make sure your .env file has the correct database credentials

Write-Host "Running database migrations..." -ForegroundColor Cyan

# Load environment variables from .env
if (Test-Path ".env") {
    Get-Content ".env" | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($key, $value)
        }
    }
    Write-Host "✓ Loaded environment variables" -ForegroundColor Green
} else {
    Write-Host "✗ .env file not found!" -ForegroundColor Red
    exit 1
}

# Get database credentials from environment
$DB_USER = $env:POSTGRES_USER
$DB_PASSWORD = $env:POSTGRES_PASSWORD
$DB_HOST = $env:POSTGRES_HOST
$DB_PORT = $env:POSTGRES_PORT
$DB_NAME = $env:POSTGRES_DBNAME

if (-not $DB_USER -or -not $DB_HOST -or -not $DB_NAME) {
    Write-Host "✗ Missing database environment variables!" -ForegroundColor Red
    Write-Host "Required: POSTGRES_USER, POSTGRES_HOST, POSTGRES_DBNAME" -ForegroundColor Yellow
    exit 1
}

Write-Host "Database: $DB_NAME at $DB_HOST:$DB_PORT" -ForegroundColor Cyan

# Set PGPASSWORD environment variable for psql
$env:PGPASSWORD = $DB_PASSWORD

# Run migrations
Write-Host "`nRunning migration: 001_create_users_table.sql" -ForegroundColor Yellow

try {
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "migrations/001_create_users_table.sql"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n✓ Migration completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "`n✗ Migration failed!" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "`n✗ Error running migration: $_" -ForegroundColor Red
    Write-Host "`nMake sure PostgreSQL client (psql) is installed and in your PATH" -ForegroundColor Yellow
    exit 1
}

# Clear password from environment
$env:PGPASSWORD = $null

Write-Host "`n✓ All done!" -ForegroundColor Green
