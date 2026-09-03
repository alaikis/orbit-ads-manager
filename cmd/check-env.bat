@echo off
chcp 65001 >nul
echo ==========================================
echo   Orbit 广告智能中枢 - 环境检查与安装
echo ==========================================
echo.

:: Check Administrator
net session >nul 2>&1
if %errorLevel% == 0 (
    echo [OK] 管理员权限
) else (
    echo [WARN] 建议以管理员权限运行此脚本
)

:: Check Go
go version >nul 2>&1
if %errorLevel% == 0 (
    for /f "tokens=3" %%g in ('go version') do echo [OK] Go %%g
) else (
    echo [WARN] Go 未安装，请从 https://go.dev/dl/ 下载 Go 1.22+
)

:: Check Node.js
node --version >nul 2>&1
if %errorLevel% == 0 (
    for /f "tokens=1" %%n in ('node --version') do echo [OK] Node.js %%n
) else (
    echo [WARN] Node.js 未安装，请从 https://nodejs.org/ 下载 Node.js 18+
)

:: Check npm
npm --version >nul 2>&1
if %errorLevel% == 0 (
    for /f "tokens=*" %%v in ('npm --version') do echo [OK] NPM %%v
) else (
    echo [WARN] NPM 未安装
)

:: Check Git
git --version >nul 2>&1
if %errorLevel% == 0 (
    for /f "tokens=*" %%v in ('git --version') do echo [OK] %%v
) else (
    echo [WARN] Git 未安装，请从 https://git-scm.com/ 下载安装
)

:: Check Docker
docker --version >nul 2>&1
if %errorLevel% == 0 (
    for /f "tokens=*" %%v in ('docker --version') do echo [OK] %%v
) else (
    echo [WARN] Docker 未安装，请从 https://www.docker.com/ 下载安装 Docker Desktop
)

:: Check PostgreSQL
docker ps --filter "name=postgres" --format "{{.Names}}" 2>nul | findstr postgres >nul
if %errorLevel% == 0 (
    echo [OK] PostgreSQL Docker 容器运行中
) else (
    echo [INFO] PostgreSQL 未运行，可使用: docker run --name orbit-pg -e POSTGRES_PASSWORD=adm2026@ -p 5432:5432 -d postgres:15-alpine
)

:: Check Redis
docker ps --filter "name=redis" --format "{{.Names}}" 2>nul | findstr redis >nul
if %errorLevel% == 0 (
    echo [OK] Redis Docker 容器运行中
) else (
    echo [INFO] Redis 未运行，可使用: docker run --name orbit-redis -p 6379:6379 -d redis:7-alpine
)

echo.
echo ==========================================
echo   环境检查完成
echo ==========================================
pause
