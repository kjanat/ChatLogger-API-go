@echo off
REM version_update.bat - A script to update the version in internal/version/version.go

if "%~1"=="" (
  echo Usage: version_update.bat ^<new_version^>
  echo Example: version_update.bat 1.2.3
  exit /b 1
)

set NEW_VERSION=%~1

REM Validate version format using simple pattern matching (not as robust as regex)
echo %NEW_VERSION% | findstr /r "^[0-9]*\.[0-9]*\.[0-9]*$" >nul
if errorlevel 1 (
  echo Error: Version must follow semantic versioning format (e.g., 1.2.3)
  exit /b 1
)

REM Get the directory of the script and project root
set "SCRIPT_DIR=%~dp0"
set "ROOT_DIR=%SCRIPT_DIR%..\"

REM Update version.go
echo Updating version in internal/version/version.go...
powershell -Command "(Get-Content '%ROOT_DIR%internal\version\version.go') -replace 'Version = \"[0-9]*\.[0-9]*\.[0-9]*\"', 'Version = \"%NEW_VERSION%\"' | Set-Content '%ROOT_DIR%internal\version\version.go'"

echo Version updated to %NEW_VERSION%
echo Remember to commit these changes and push them to your repository
echo You can also run the tag-version GitHub Action workflow to create a git tag
