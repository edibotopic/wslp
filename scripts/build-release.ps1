$root = Split-Path $PSScriptRoot -Parent
$releaseDir = Join-Path $root "release"

# Clean and recreate release dir
if (Test-Path $releaseDir) {
    Remove-Item $releaseDir -Recurse -Force
}
New-Item -ItemType Directory -Path $releaseDir | Out-Null

# Build CLI
Write-Host "Building wslp.exe..."
Push-Location $root
go build -o "$releaseDir\wslp.exe"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to build wslp.exe"
    exit 1
}
Write-Host "  wslp.exe built successfully"

# Build GUI
Write-Host "Building GUI..."
Push-Location (Join-Path $root "gui")
flutter build windows --release
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to build GUI"
    Pop-Location
    exit 1
}
Pop-Location
Write-Host "  GUI built successfully"

# Zip GUI + rungui.bat + wslp.exe into wslp-full.zip
Write-Host "Creating wslp-full.zip..."
$flutterRelease = Join-Path $root "gui\build\windows\x64\runner\Release"
$zipStaging = Join-Path $root "release\_staging"
New-Item -ItemType Directory -Path $zipStaging | Out-Null
Copy-Item -Path "$flutterRelease\*" -Destination $zipStaging -Recurse
Copy-Item -Path (Join-Path $PSScriptRoot "rungui-release.bat") -Destination "$zipStaging\rungui.bat"
Copy-Item -Path "$releaseDir\wslp.exe" -Destination $zipStaging
Compress-Archive -Path "$zipStaging\*" -DestinationPath "$releaseDir\wslp-full.zip"
Remove-Item $zipStaging -Recurse -Force
Pop-Location

Write-Host ""
Write-Host "Release assets ready in $($releaseDir):"
Write-Host "  wslp.exe       -- CLI only"
Write-Host "  wslp-full.zip  -- CLI + GUI (extract and run rungui.bat)"
