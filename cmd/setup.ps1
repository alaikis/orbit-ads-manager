# Orbit Environment Setup Script for Windows
# 自动检测并安装 Go, Node.js, PostgreSQL, Redis

param(
    [switch]$SkipGo,
    [switch]$SkipNode,
    [switch]$SkipPostgres,
    [switch]$SkipRedis,
    [switch]$Help
)

$ErrorActionPreference = "Stop"

function Show-Help {
    Write-Host @"
Orbit 环境自动安装脚本

用法:
  .\cmd\setup.ps1 [选项]

选项:
  -SkipGo        跳过 Go 安装
  -SkipNode      跳过 Node.js 安装
  -SkipPostgres  跳过 PostgreSQL 安装
  -SkipRedis     跳过 Redis 安装
  -Help          显示此帮助信息

说明:
  此脚本会自动下载并安装所需的开发环境。
  需要管理员权限运行。
"@
    exit 0
}

if ($Help) { Show-Help }

# Check for administrator privileges
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Warn "建议以管理员权限运行此脚本"
}

$ProgressPreference = "SilentlyContinue"

function Test-Command {
    param([string]$cmd)
    try {
        Get-Command $cmd -ErrorAction SilentlyContinue | Select-Object -First 1
        return $true
    } catch {
        return $false
    }
}

function Get-InstalledVersion {
    param([string]$cmd, [string]$pattern)
    if (Test-Command $cmd) {
        $output = & $cmd --version 2>&1
        if ($output -match $pattern) { return $Matches[1] }
    }
    return $null
}

# Check Go
Write-Host ""
Write-Host "检查 Go..." -ForegroundColor Cyan
$goVersion = Get-InstalledVersion "go" "go(\d+\.\d+)"
if ($goVersion) {
    Write-Ok "Go $goVersion 已安装"
} elseif (-not $SkipGo) {
    Write-Info "Go 未安装，正在下载..."
    $goUrl = "https://go.dev/dl/go1.22.5.windows-amd64.msi"
    $goInstaller = "$env:TEMP\go-installer.msi"
    Invoke-WebRequest -Uri $goUrl -OutFile $goInstaller
    Write-Info "安装 Go..."
    Start-Process msiexec.exe -Wait -ArgumentList "/i `"$goInstaller`" /quiet /norestart"
    Remove-Item $goInstaller
    Write-Ok "Go 安装完成，请重启终端后运行 go version 验证"
} else {
    Write-Warn "跳过 Go 安装"
}

# Check Node.js
Write-Host ""
Write-Host "检查 Node.js..." -ForegroundColor Cyan
$nodeVersion = (node --version 2>&1) -replace 'v', ''
if ($nodeVersion -match '^\d+') {
    $major = [int]$Matches[0]
    if ($major -ge 18) {
        Write-Ok "Node.js $nodeVersion 已安装"
    } else {
        Write-Warn "Node.js 版本过低 ($nodeVersion)，需要 18+"
        if (-not $SkipNode) {
            Write-Info "请从 https://nodejs.org/ 下载安装 Node.js 18+ LTS"
        }
    }
} else {
    Write-Warn "Node.js 未安装"
    if (-not $SkipNode) {
        Write-Info "请从 https://nodejs.org/ 下载安装 Node.js 18+ LTS"
    }
}

# Check PostgreSQL
Write-Host ""
Write-Host "检查 PostgreSQL..." -ForegroundColor Cyan
if (Get-Service -Name "postgresql*" -ErrorAction SilentlyContinue) {
    Write-Ok "PostgreSQL 服务已安装"
} else {
    Write-Warn "PostgreSQL 未安装"
    if (-not $SkipPostgres) {
        Write-Info "请从 https://www.postgresql.org/download/ 下载安装 PostgreSQL 15+"
        Write-Info "或使用 Docker: docker run --name orbit-pg -e POSTGRES_PASSWORD=adm2026@ -p 5432:5432 -d postgres:15-alpine"
    }
}

# Check Redis
Write-Host ""
Write-Host "检查 Redis..." -ForegroundColor Cyan
if (Get-Service -Name "redis*" -ErrorAction SilentlyContinue) {
    Write-Ok "Redis 服务已安装"
} else {
    Write-Warn "Redis 未安装"
    if (-not $SkipRedis) {
        Write-Info "请从 https://redis.io/download/ 下载安装 Redis"
        Write-Info "或使用 Docker: docker run --name orbit-redis -p 6379:6379 -d redis:7-alpine"
    }
}

# Check Git
Write-Host ""
Write-Host "检查 Git..." -ForegroundColor Cyan
if (Test-Command "git") {
    $gitVersion = git --version
    Write-Ok "$gitVersion 已安装"
} else {
    Write-Warn "Git 未安装，请从 https://git-scm.com/ 下载安装"
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  环境检查完成" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
