@echo off
REM Simple batch file to run integration tests for LUGX Gaming services

echo LUGX Gaming - Integration Test Runner
echo =====================================

set "SCRIPT_DIR=%~dp0"

if "%1"=="game" goto :run_game
if "%1"=="order" goto :run_order  
if "%1"=="analytics" goto :run_analytics
if "%1"=="help" goto :show_help
if "%1"=="/?" goto :show_help
if "%1"=="-h" goto :show_help
if "%1"=="--help" goto :show_help

echo Running all integration tests...
echo.

:run_game
echo Testing Game Service...
cd /d "%SCRIPT_DIR%game-service"
go test -v
if errorlevel 1 (
    echo Game Service tests failed!
    set "TESTS_FAILED=1"
) else (
    echo Game Service tests passed!
)
echo.

if "%1"=="game" goto :end

:run_order
echo Testing Order Service...
cd /d "%SCRIPT_DIR%order-service"
go test -v
if errorlevel 1 (
    echo Order Service tests failed!
    set "TESTS_FAILED=1"
) else (
    echo Order Service tests passed!
)
echo.

if "%1"=="order" goto :end

:run_analytics
echo Testing Analytics Service...
cd /d "%SCRIPT_DIR%analytics-service"
go test -v
if errorlevel 1 (
    echo Analytics Service tests failed!
    set "TESTS_FAILED=1"
) else (
    echo Analytics Service tests passed!
)
echo.

if "%1"=="analytics" goto :end

goto :end

:show_help
echo Usage: run-tests.bat [service]
echo.
echo Arguments:
echo   game        Run only Game Service tests
echo   order       Run only Order Service tests  
echo   analytics   Run only Analytics Service tests
echo   (no args)   Run all tests
echo   help        Show this help message
echo.
echo Examples:
echo   run-tests.bat           - Run all tests
echo   run-tests.bat game      - Run only game service tests
echo   run-tests.bat help      - Show this help
echo.
goto :end

:end
cd /d "%SCRIPT_DIR%"

if defined TESTS_FAILED (
    echo.
    echo Some tests failed. Check the output above for details.
    exit /b 1
) else (
    echo.
    echo All tests completed successfully!
    exit /b 0
)
