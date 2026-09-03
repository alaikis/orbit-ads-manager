# Orbit Build Script for Windows (PowerShell 5.1+)
# Auto detect env, install deps, build frontend/backend, output to release/

param(
    [switch]$SkipFrontend,
    [switch]$SkipBackend,
    [switch]$SkipDeps,
    [switch]$Clean,
    [switch]$Help
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir
$ReleaseDir = Join-Path $RootDir "release"

function Write-Info { param([string]$msg) Write-Host "[INFO] $msg" -ForegroundColor Cyan }
function Write-Ok { param([string]$msg) Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-Warn { param([string]$msg) Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Write-Err { param([string]$msg) Write-Host "[ERROR] $msg" -ForegroundColor Red }

function Show-Help {
    Write-Host @"
Orbit Automated Build Script

Usage:
  .\cmd\build.ps1 [Options]

Options:
  -SkipFrontend   Skip frontend build
  -SkipBackend    Skip backend build
  -SkipDeps       Skip dependency installation
  -Clean          Clean build artifacts before building
  -Help           Show this help message

Examples:
  .\cmd\build.ps1                # Full build (deps + frontend + backend)
  .\cmd\build.ps1 -Clean         # Clean and full build
  .\cmd\build.ps1 -SkipFrontend  # Backend only
"@
}

function Get-BuildInfo {
    $os = "windows"
    $arch = "amd64"
    if ([Environment]::Is64BitOperatingSystem) {
        $env:PROCESSOR_ARCHITECTURE = $env:PROCESSOR_ARCHITECTURE.ToLower()
        if ($env:PROCESSOR_ARCHITECTURE -like "*arm*") {
            $arch = "arm64"
        }
    }
    $goVersion = go version 2>$null | Select-String -Pattern "go(\d+\.\d+)" | ForEach-Object { $_.Matches[0].Groups[1].Value }
    $nodeVersion = node --version 2>$null
    $npmVersion = npm --version 2>$null

    return [PSCustomObject]@{
        OS       = $os
        Arch     = $arch
        Go       = $goVersion
        Node     = $nodeVersion
        NPM      = $npmVersion
    }
}

function Test-Environment {
    param([PSCustomObject]$BuildInfo)

    $missing = @()

    if (-not $BuildInfo.Go) {
        $missing += "Go 1.22+ (https://go.dev/dl/)"
    } else {
        $goMajorMinor = $BuildInfo.Go.Split('.')[0..1] -join '.'
        if ([version]$goMajorMinor -lt [version]"1.22") {
            $missing += "Go $($BuildInfo.Go) is too old, need Go 1.22+"
        }
    }

    if (-not $BuildInfo.Node) {
        $missing += "Node.js 18+ (https://nodejs.org/)"
    } else {
        $nodeMajor = [int]($BuildInfo.Node.TrimStart('v').Split('.')[0])
        if ($nodeMajor -lt 18) {
            $missing += "Node.js $($BuildInfo.Node) is too old, need Node.js 18+"
        }
    }

    $pgRunning = Get-Process -Name "postgres" -ErrorAction SilentlyContinue
    if (-not $pgRunning) {
        Write-Warn "PostgreSQL not detected, build will continue but local testing needs database"
    }

    $redisRunning = Get-Process -Name "redis-server" -ErrorAction SilentlyContinue
    if (-not $redisRunning) {
        Write-Warn "Redis not detected, build will continue but local testing needs Redis"
    }

    if ($missing.Count -gt 0) {
        Write-Err "Missing dependencies:"
        $missing | ForEach-Object { Write-Host "  - $_" }
        Write-Host ""
        Write-Host "Install missing dependencies and retry, or use -SkipDeps to skip check" -ForegroundColor Yellow
        exit 1
    }

    Write-Ok "Environment check passed"
    Write-Info "  OS: $($BuildInfo.OS) | Arch: $($BuildInfo.Arch) | Go: $($BuildInfo.Go) | Node: $($BuildInfo.Node) | NPM: $($BuildInfo.NPM)"
}

function Invoke-Clean {
    param([string]$RootDir, [string]$ReleaseDir)

    Write-Info "Cleaning build artifacts..."

    $pathsToClean = @(
        (Join-Path $RootDir "apps\api\bin"),
        (Join-Path $RootDir "apps\web\.next"),
        (Join-Path $RootDir "apps\web\out"),
        (Join-Path $RootDir "apps\web\dist"),
        $ReleaseDir
    )

    foreach ($path in $pathsToClean) {
        if (Test-Path $path) {
            Remove-Item -Recurse -Force $path
            Write-Info "  Removed: $path"
        }
    }

    Write-Ok "Clean complete"
}

function Invoke-BackendDeps {
    param([string]$ApiDir)

    Write-Info "Installing backend dependencies..."
    Set-Location $ApiDir
    go mod tidy
    Write-Ok "Backend dependencies installed"
}

function Invoke-FrontendDeps {
    param([string]$WebDir)

    Write-Info "Installing frontend dependencies..."
    Set-Location $WebDir
    if (-not (Test-Path "node_modules")) {
        npm install
    } else {
        Write-Info "  node_modules exists, skipping install"
    }
    Write-Ok "Frontend dependencies installed"
}

function Invoke-BackendBuild {
    param([string]$ApiDir, [string]$ReleaseDir, [string]$OS, [string]$Arch)

    Write-Info "Building unified backend binary (API + Worker + Scheduler)..."

    $binDir = Join-Path $ReleaseDir "bin"
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null

    Set-Location $ApiDir

    Write-Info "  Building server..."
    $serverOutput = Join-Path $binDir "orbit-server-$OS-$Arch.exe"
    go build -ldflags "-s -w" -o $serverOutput ./cmd/server
    if ($LASTEXITCODE -ne 0) {
        Write-Err "server build failed"
        exit 1
    }
    Write-Ok "  server -> $serverOutput"

    Copy-Item (Join-Path $ApiDir ".env.example") (Join-Path $ReleaseDir ".env.example") -Force

    Write-Ok "Backend build complete"
}

function Invoke-FrontendBuild {
    param([string]$WebDir, [string]$ReleaseDir)

    Write-Info "Building frontend..."

    Set-Location $WebDir

    Write-Info "  Running lint..."
    npm run lint 2>&1 | Out-Null

    Write-Info "  Running typecheck..."
    npm run typecheck 2>&1 | Out-Null

    Write-Info "  Running build..."
    npm run build 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Err "Frontend build failed"
        exit 1
    }

    $frontendOutput = Join-Path $ReleaseDir "web"
    if (Test-Path $frontendOutput) {
        Remove-Item -Recurse -Force $frontendOutput
    }

    if (Test-Path (Join-Path $WebDir ".next")) {
        Copy-Item (Join-Path $WebDir ".next") $frontendOutput -Recurse -Force
    }
    if (Test-Path (Join-Path $WebDir "public")) {
        Copy-Item (Join-Path $WebDir "public") (Join-Path $frontendOutput "public") -Recurse -Force
    }
    Copy-Item (Join-Path $WebDir "package.json") (Join-Path $frontendOutput "package.json") -Force

    Write-Ok "Frontend build complete -> $frontendOutput"
}

function New-DeploymentScripts {
    param([string]$ReleaseDir, [string]$OS)

    Write-Info "Generating deployment scripts..."

    $compose = @"
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: adm
      POSTGRES_PASSWORD: adm2026@
      POSTGRES_DB: adshub
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U adm"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    image: orbit-api:latest
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://adm:adm2026%40portal.gusty.top:5432/adshub
      - REDIS_URL=redis://redis:6379/0
      - APP_BASE_URL=https://orbit.elstella.com
      - API_BASE_URL=https://api.orbit.elstella.com
      - PORT=8080
      - JWT_SECRET=`$JWT_SECRET
      - ENCRYPTION_MASTER_KEY=`$ENCRYPTION_MASTER_KEY
      - TZ=Asia/Shanghai
      - LOG_LEVEL=info
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health/ready"]
      interval: 10s
      timeout: 5s
      retries: 3

  worker:
    image: orbit-worker:latest
    restart: unless-stopped
    environment:
      - DATABASE_URL=postgres://adm:adm2026%40portal.gusty.top:5432/adshub
      - REDIS_URL=redis://redis:6379/0
      - JWT_SECRET=`$JWT_SECRET
      - ENCRYPTION_MASTER_KEY=`$ENCRYPTION_MASTER_KEY
      - TZ=Asia/Shanghai
      - LOG_LEVEL=info
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  scheduler:
    image: orbit-scheduler:latest
    restart: unless-stopped
    environment:
      - DATABASE_URL=postgres://adm:adm2026%40portal.gusty.top:5432/adshub
      - REDIS_URL=redis://redis:6379/0
      - JWT_SECRET=`$JWT_SECRET
      - ENCRYPTION_MASTER_KEY=`$ENCRYPTION_MASTER_KEY
      - TZ=Asia/Shanghai
      - LOG_LEVEL=info
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

volumes:
  postgres_data:
  redis_data:
"@

    $compose | Out-File -FilePath (Join-Path $ReleaseDir "docker-compose.yml") -Encoding UTF8

    $startScript = @"
@echo off
echo ==========================================
echo   Orbit Ad Intelligence Platform
echo ==========================================
echo.

if not exist ".env" (
    echo [ERROR] .env file missing
    pause
    exit /b 1
)

echo [INFO] Starting services...
docker compose up -d

echo.
echo [OK] Services started
echo.
echo URLs:
echo   - Frontend: http://localhost:3000
echo   - Backend: http://localhost:8080/api/v1
echo   - Health: http://localhost:8080/health/ready
echo.
pause
"@

    $startScript | Out-File -FilePath (Join-Path $ReleaseDir "start.bat") -Encoding ASCII

    $stopScript = @"
@echo off
echo [INFO] Stopping services...
docker compose down
echo [OK] Services stopped
pause
"@

    $stopScript | Out-File -FilePath (Join-Path $ReleaseDir "stop.bat") -Encoding ASCII

    Write-Ok "Deployment scripts generated"
}

function Main {
    if ($Help) {
        Show-Help
        exit 0
    }

    Write-Host ""
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "  Orbit Ad Intelligence Platform - Build" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host ""

    $buildInfo = Get-BuildInfo

    if (-not $SkipDeps) {
        Test-Environment -BuildInfo $buildInfo
    }

    if ($Clean) {
        Invoke-Clean -RootDir $RootDir -ReleaseDir $ReleaseDir
    }

    if (-not $SkipBackend) {
        Invoke-BackendDeps -ApiDir (Join-Path $RootDir "apps\api")
        Invoke-BackendBuild -ApiDir (Join-Path $RootDir "apps\api") -ReleaseDir $ReleaseDir -OS $buildInfo.OS -Arch $buildInfo.Arch
    }

    if (-not $SkipFrontend) {
        Invoke-FrontendDeps -WebDir (Join-Path $RootDir "apps\web")
        Invoke-FrontendBuild -WebDir (Join-Path $RootDir "apps\web") -ReleaseDir $ReleaseDir
    }

    New-DeploymentScripts -ReleaseDir $ReleaseDir -OS $buildInfo.OS

    Write-Host ""
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host "  Build Complete!" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host ""
    Write-Info "Output: $ReleaseDir"
    Write-Info "Contents:"
    Write-Info "  - bin/orbit-server-*.exe    (Unified Server: API + Worker + Scheduler)"
    Write-Info "  - web/                     (Frontend Static Files)"
    Write-Info "  - docker-compose.yml       (Production Deploy)"
    Write-Info "  - start.bat / stop.bat     (Quick Start/Stop)"
    Write-Host ""
    Write-Info "Deploy Steps:"
    Write-Info "  1. Upload release/ directory to server"
    Write-Info "  2. Copy .env.example to .env and configure environment variables"
    Write-Info "  3. Run start.bat or docker compose up -d"
    Write-Host ""
}

Main
