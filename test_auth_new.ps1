# PowerShell script to test authentication endpoints
# Run this after starting your server

$BASE_URL = "http://localhost:8080"

Write-Host "Testing TetriON Authentication Endpoints" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Health Check
Write-Host "1. Testing Health Check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/api/health" -Method Get
    Write-Host "OK Health Check: $($response.status)" -ForegroundColor Green
} catch {
    Write-Host "X Health Check Failed: $_" -ForegroundColor Red
}
Write-Host ""

# Test 2: Register User
Write-Host "2. Testing User Registration..." -ForegroundColor Yellow
$registerBody = @{
    username = "testuser"
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/api/auth/register" -Method Post -Body $registerBody -ContentType "application/json"
    Write-Host "OK Registration successful!" -ForegroundColor Green
    Write-Host "  User ID: $($response.user.id)" -ForegroundColor Gray
    Write-Host "  Username: $($response.user.username)" -ForegroundColor Gray
    $token = $response.token
    Write-Host "  Token: $($token.Substring(0, 20))..." -ForegroundColor Gray
} catch {
    Write-Host "X Registration Failed: $_" -ForegroundColor Red
    $response = $null
    $token = $null
}
Write-Host ""

# Test 3: Login
Write-Host "3. Testing User Login..." -ForegroundColor Yellow
$loginBody = @{
    username = "testuser"
    password = "password123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/api/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    Write-Host "OK Login successful!" -ForegroundColor Green
    Write-Host "  Username: $($response.user.username)" -ForegroundColor Gray
    $token = $response.token
    Write-Host "  Token: $($token.Substring(0, 20))..." -ForegroundColor Gray
} catch {
    Write-Host "X Login Failed: $_" -ForegroundColor Red
}
Write-Host ""

# Test 4: Get Profile (Protected Endpoint)
if ($token) {
    Write-Host "4. Testing Protected Profile Endpoint..." -ForegroundColor Yellow
    try {
        $headers = @{
            "Authorization" = "Bearer $token"
        }
        $response = Invoke-RestMethod -Uri "$BASE_URL/api/auth/profile" -Method Get -Headers $headers
        Write-Host "OK Profile retrieved successfully!" -ForegroundColor Green
        Write-Host "  Username: $($response.user.username)" -ForegroundColor Gray
        Write-Host "  Email: $($response.user.email)" -ForegroundColor Gray
    } catch {
        Write-Host "X Profile Request Failed: $_" -ForegroundColor Red
    }
} else {
    Write-Host "4. Skipping profile test (no token available)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Testing Complete!" -ForegroundColor Cyan
