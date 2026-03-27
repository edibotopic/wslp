@echo off
echo Starting WSL Plus GUI...

REM Start the server in the background
start /B wslp.exe serve

REM Wait a moment for server to start
timeout /t 2 /nobreak >nul

REM Launch the GUI (same folder as this script)
start "" "%~dp0gui.exe"

echo WSL Plus is running.
echo Close this window to stop the server.
echo.

REM Keep the server running until user closes this window
pause >nul

REM Kill the server process when done
taskkill /F /IM wslp.exe >nul 2>&1
