# version_update.ps1 - A PowerShell script to update the version in internal/version/version.go

param (
    [Parameter(Mandatory, ValueFromRemainingArguments)]
    [string]$NewVersion
)

# Validate version format (Semantic Versioning: major.minor.patch)
if ($NewVersion -notmatch '^\d+\.\d+\.\d+$') {
    Write-Error "Error: Version must follow semantic versioning format (e.g., 1.2.3)"
    exit 1
}

# Get the directory of the script and set the project root
$ScriptDir = Split-Path -Path $MyInvocation.MyCommand.Definition -Parent
$RootDir = Split-Path -Path $ScriptDir -Parent

# Version file path
$VersionFilePath = Join-Path -Path $RootDir -ChildPath "internal\version\version.go"

Write-Host "Updating version in internal/version/version.go..."

# Read version.go content
$Content = Get-Content -Path $VersionFilePath -Raw

# Replace version string
$UpdatedContent = $Content -replace 'Version = "[0-9]+\.[0-9]+\.[0-9]+"', "Version = `"$NewVersion`""

# Write updated content back to file
$UpdatedContent | Set-Content -Path $VersionFilePath -NoNewline

Write-Host "Version updated to $NewVersion"
Write-Host "Remember to commit these changes and push them to your repository"
Write-Host "You can also run the tag-version GitHub Action workflow to create a git tag"
