@echo off
echo ========================================
echo   Starting CollectBFMSRQ Application
echo ========================================
echo.

:: Start backend in a new window
echo Starting backend server...
start "Backend Server" cmd /k "cd backend && go run cmd\server\main.go"
timeout /t 2 >nul

:: Start frontend in a new window
echo Starting frontend server...
start "Frontend Server" cmd /k "cd frontend && npm run dev"
timeout /t 2 >nul

echo.
echo ========================================
echo   Application Started!
echo ========================================
echo.
echo Backend: http://localhost:8080
echo Frontend: http://localhost:3000
echo.
echo Both servers are running in separate windows.
echo Close those windows to stop the servers.
echo.
pause
