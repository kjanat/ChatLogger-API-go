#!/bin/pwsh
param (
	[Parameter(HelpMessage = "Generate OpenAPI v3.1 documentation, defaults to v2 which is better supported by swag.")]
	[switch] $v31,
	[Parameter(HelpMessage = "Output path for the generated documentation.")]
	[string] $outPath = "./docs",
	[Parameter(HelpMessage = "Instance name for the generated documentation.")]
	[ValidateSet('OpenAPI', 'OpenAPIv3')]
	[string] $instanceName = 'OpenAPI'
	<# ,
	[Parameter(HelpMessage = "Package name for the generated documentation.")]
	[ValidateSet('OpenAPI', 'OpenAPI_v3')]
	[string] $packageName = 'OpenAPI'
	#>
)
# This script generates OpenAPI 3.1 documentation from Go annotations using swag v2

# Ensure we're starting with fresh docs
if (Test-Path "$PSScriptRoot/../$outPath") {
    Write-Host "Cleaning previous documentation..."
    Get-ChildItem -Recurse -Path "$PSScriptRoot/../$outPath/${instanceName}_*" | Remove-Item -Force
}

# Ensure the output directory exists
New-Item -ItemType Directory -Path "$PSScriptRoot/../$outPath" -Force | Out-Null

# Display current swag version
$swagVersion = $(swag --version)
Write-Host "Swag version: " $swagVersion

# Change directory to project root to ensure proper parsing
Push-Location "$PSScriptRoot/.."
try {
	switch ($v31) {
		$true {
			Write-Host "Generating OpenAPI v3 documentation..."
			swag init `
				--output $outPath `
				<# --generalInfo 'cmd/server/main.go' #> `
				--dir 'cmd/server,internal/api,internal/config,internal/domain,internal/handler,internal/hash,internal/jobs,internal/middleware,internal/repository,internal/service,internal/strategy,internal/version' `
				--exclude $outPath `
				--instanceName $instanceName --v3.1 --packageName "docs_v3" `
				--parseInternal --parseFuncBody --requiredByDefault `
				#--state value                             # Set host state for swagger.json
		}
		Default {
			Write-Host "Generating OpenAPI v2 documentation..."
			swag init `
				--output $outPath `
				<# --generalInfo 'cmd/server/main.go' #> `
				--dir 'cmd/server,internal/api,internal/config,internal/domain,internal/handler,internal/hash,internal/jobs,internal/middleware,internal/repository,internal/service,internal/strategy,internal/version' `
				--exclude $outPath `
				--instanceName $instanceName --packageName "docs" `
				--parseInternal --parseFuncBody --requiredByDefault `
				#--state value                             # Set host state for swagger.json
		}
	}

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Documentation generation failed with code $LASTEXITCODE"
        exit $LASTEXITCODE
    }

    # Verify the documentation files were created
    if (Test-Path "${outPath}/${instanceName}_swagger.json") {
        Write-Host "✅ Documentation generated ${instanceName} successfully in ${outPath}/" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to generate documentation files" -ForegroundColor Red
        exit 1
    }
} finally {
    # Return to original directory
    Pop-Location
}
