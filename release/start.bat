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
