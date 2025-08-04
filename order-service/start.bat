@echo off
echo Starting LUGX Gaming - Order Service...
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Check if PostgreSQL is running (optional check)
echo Checking for PostgreSQL...
netstat -an | find "5432" >nul
if errorlevel 1 (
    echo Warning: PostgreSQL might not be running on port 5432
    echo Make sure PostgreSQL is installed and running, or use Docker Compose
    echo.
)

REM Install dependencies
echo Installing Go dependencies...
go mod tidy
if errorlevel 1 (
    echo Error: Failed to install dependencies
    pause
    exit /b 1
)

echo.
echo Dependencies installed successfully!
echo.

REM Start the service
echo Starting Order Service on port 8081...
echo Press Ctrl+C to stop the service
echo.
go run main.go

pause
