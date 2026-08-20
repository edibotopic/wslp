@echo off
REM Optional: pass a port number to run on something other than 8080,
REM e.g. `rungui.bat 9090`.
set PORT=%1
if "%PORT%"=="" set PORT=8080

echo Starting WSL Plus GUI on port %PORT%...

REM Start the server in the background
start /B wslp.exe serve --port %PORT%

REM Wait a moment for server to start
timeout /t 2 /nobreak >nul

REM Launch the GUI, pointing it at the same port. The server also stops
REM itself automatically when the GUI window is closed, so this is mainly
REM a fallback for cases where the console window is closed instead.
start "" "gui\build\windows\x64\runner\Release\gui.exe" --port=%PORT%

echo WSL Plus is running.
echo Close this window to stop the server.
echo.

REM Keep the server running until user closes this window
pause >nul

REM Kill the server process when done
taskkill /F /IM wslp.exe >nul 2>&1
