@echo off
echo Building WSL Plus...
echo.

echo [1/2] Building wslp.exe (Go backend)...
go build -o wslp.exe
if %errorlevel% neq 0 (
    echo ERROR: Failed to build wslp.exe
    exit /b 1
)
echo ✓ wslp.exe built successfully
echo.

echo [2/2] Building GUI (Flutter)...
cd gui
call flutter build windows --release
if %errorlevel% neq 0 (
    echo ERROR: Failed to build GUI
    cd ..
    exit /b 1
)
cd ..
echo ✓ GUI built successfully
echo.

echo ============================================
echo Build complete!
echo.
echo To use WSL Plus:
echo   - CLI: Run 'wslp.exe' from the command line
echo   - GUI: Run 'rungui.bat'
echo ============================================
pause
