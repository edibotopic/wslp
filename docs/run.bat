@echo off
setlocal

set SPHINXDIR=.sphinx
set VENVDIR=%SPHINXDIR%\venv
set BUILDDIR=_build
set SPHINX_HOST=127.0.0.1
set SPHINX_PORT=8000

echo [1/3] Generating CLI reference documentation...
pushd ..
go run ./cmd/docgen -out ./docs/reference -format markdown
if %errorlevel% neq 0 (
    echo ERROR: Failed to generate CLI reference docs
    popd
    exit /b 1
)
popd
echo CLI reference generated.
echo.

echo [2/3] Setting up Python virtualenv...
if not exist "%VENVDIR%\Scripts\sphinx-autobuild.exe" (
    echo Creating virtualenv and installing dependencies...
    if exist "%VENVDIR%" rmdir /s /q "%VENVDIR%"
    python -m venv "%VENVDIR%"
    if %errorlevel% neq 0 (
        echo ERROR: Failed to create virtualenv. Make sure Python is installed.
        exit /b 1
    )
    "%VENVDIR%\Scripts\pip.exe" install --upgrade -r requirements.txt
    if %errorlevel% neq 0 (
        echo ERROR: Failed to install requirements.
        exit /b 1
    )
)
echo Virtualenv ready.
echo.

echo [3/3] Starting sphinx-autobuild at http://%SPHINX_HOST%:%SPHINX_PORT%
"%VENVDIR%\Scripts\sphinx-autobuild.exe" -b dirhtml --host %SPHINX_HOST% --port %SPHINX_PORT% "." "%BUILDDIR%" -c . -d "%SPHINXDIR%\.doctrees" -j auto
