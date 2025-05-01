#!/bin/pwsh
# scripts/docs_generate.ps1
# This script generates Swagger/OpenAPI documentation from Go annotations

# Output location for the generated docs
$outPath = "cmd/docs"

# Ensure the output directory exists
if (-not (Test-Path "$PSScriptRoot/../$outPath")) {
    New-Item -ItemType Directory -Path "$PSScriptRoot/../$outPath" -Force | Out-Null
}

# Display current swag version
$swagVersion = $(swag --version)
Write-Host "Swag version: " $swagVersion

# Generate OpenAPI 3.1 docs
# NOTE: No need to set version manually - it will be set at runtime from version.Version
$arguments = @(
    "--v3.1",                                  # Generate OpenAPI 3.1
    "--generalInfo", "./cmd/server/main.go",   # Where to find general API info
    "--output", $outPath,                     # Where to output docs
    "--parseDependency"                        # Parse code in dependency directories
)

Write-Host "Generating API documentation..."
swag init @arguments

Write-Host "API documentation generated in $outPath/"
