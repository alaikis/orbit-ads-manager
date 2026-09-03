@echo off
chcp 65001 >nul
echo ==========================================
echo   Orbit 广告智能中枢 - 本地开发启动
echo ==========================================
echo.

:: Check if Docker is available
docker --version >nul 2>&1
if %errorLevel% == 0 (
    echo [INFO] 检测到 Docker，启动 PostgreSQL 和 Redis...
    docker compose up -d postgres redis
    echo [OK] 数据库服务已启动
    echo.
)

:: Start backend
echo [INFO] 启动后端服务...
start "Orbit Backend" cmd /c "cd apps/api && go run cmd/server/main.go"

:: Wait a moment for backend to start
timeout /t 3 /nobreak >nul

:: Start frontend
echo [INFO] 启动前端服务...
start "Orbit Frontend" cmd /c "cd apps/web && npm run dev"

echo.
echo ==========================================
echo   服务启动完成
echo ==========================================
echo.
echo 访问地址:
echo   - 前端: http://localhost:3000
echo   - 后端: http://localhost:8080/api/v1
echo   - 健康检查: http://localhost:8080/health/ready
echo.
pause
