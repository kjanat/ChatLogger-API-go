#!/bin/pwsh
# This script generates OpenAPI 3.1 documentation from Go annotations using swag v2

# Output location - keeping it in the root "docs" directory for easier imports
$outPath = "docs"

# Ensure we're starting with fresh docs
if (Test-Path "$PSScriptRoot/../$outPath") {
    Write-Host "Cleaning previous documentation..."
    Remove-Item -Recurse -Path "$PSScriptRoot/../$outPath" -Force
}

# Ensure the output directory exists
New-Item -ItemType Directory -Path "$PSScriptRoot/../$outPath" -Force | Out-Null

# Display current swag version
$swagVersion = $(swag --version)
Write-Host "Swag version: " $swagVersion

# Change directory to project root to ensure proper parsing
Push-Location "$PSScriptRoot/.."
try {
    Write-Host "Generating OpenAPI 3.1 documentation..."
    
    # Use swag init with proper flags for API documentation
    & swag init `
        --generalInfo "./cmd/server/main.go" `
        --output "$outPath" `
        --parseVendor `
        --parseDependency `
        <# --v3.1 #> `
        --propertyStrategy "camelcase"
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Documentation generation failed with code $LASTEXITCODE"
        exit $LASTEXITCODE
    }

    # Verify the documentation files were created
    if (Test-Path "$outPath/swagger.json") {
        Write-Host "✅ Documentation generated successfully in $outPath/" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to generate documentation files" -ForegroundColor Red
        exit 1
    }
} finally {
    # Return to original directory
    Pop-Location
}
