@echo off
echo ========================================
echo   CollectBFMSRQ Setup Script
echo ========================================
echo.

:: Check if PostgreSQL is running
echo [1/5] Checking PostgreSQL...
where psql >nul 2>&1
if %errorlevel% neq 0 (
    echo WARNING: PostgreSQL not found in PATH
    echo Please install PostgreSQL and add it to your PATH
    echo Download from: https://www.postgresql.org/download/windows/
    pause
    exit /b 1
)
echo ✓ PostgreSQL found

:: Create database
echo.
echo [2/5] Setting up database...
set /p DB_USER="Enter PostgreSQL username [postgres]: "
if "%DB_USER%"=="" set DB_USER=postgres

set /p DB_PASSWORD="Enter PostgreSQL password: "

set /p DB_HOST="Enter PostgreSQL host [localhost]: "
if "%DB_HOST%"=="" set DB_HOST=localhost

set /p DB_PORT="Enter PostgreSQL port [5432]: "
if "%DB_PORT%"=="" set DB_PORT=5432

:: Create .env file
echo.
echo [3/5] Creating backend .env file...
(
echo # Server Configuration
echo SERVER_PORT=8080
echo.
echo # Database Configuration
echo DB_HOST=%DB_HOST%
echo DB_PORT=%DB_PORT%
echo DB_USER=%DB_USER%
echo DB_PASSWORD=%DB_PASSWORD%
echo DB_NAME=questionnaire_db
) > backend\.env

echo ✓ Created backend\.env

:: Install frontend dependencies
echo.
echo [4/5] Installing frontend dependencies...
cd frontend
call npm install
if %errorlevel% neq 0 (
    echo ERROR: Failed to install frontend dependencies
    pause
    exit /b 1
)
cd ..
echo ✓ Frontend dependencies installed

:: Build backend
echo.
echo [5/5] Building backend...
cd backend
go mod download
if %errorlevel% neq 0 (
    echo ERROR: Failed to download Go dependencies
    pause
    exit /b 1
)
cd ..
echo ✓ Backend dependencies downloaded

echo.
echo ========================================
echo   Setup Complete!
echo ========================================
echo.
echo Next steps:
echo 1. Create the database in PostgreSQL:
echo    psql -U %DB_USER% -h %DB_HOST% -p %DB_PORT%
echo    CREATE DATABASE questionnaire_db;
echo    \q
echo.
echo 2. Start the backend server:
echo    cd backend
echo    go run cmd\server\main.go
echo.
echo 3. In a new terminal, start the frontend:
echo    cd frontend
echo    npm run dev
echo.
echo 4. Open http://localhost:3000 in your browser
echo.
pause
