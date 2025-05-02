#!/usr/bin/env bash
set -euo pipefail

# scripts/docs_generate.sh
# This script generates Swagger/OpenAPI documentation from Go annotations

# Where this script lives
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Output location for the generated docs (relative to project root)
outPath="docs"

# Ensure the output directory exists
mkdir -p "$script_dir/../$outPath"

# Display current swag version
swagVersion="$(swag --version)"
echo "Swag version: $swagVersion"

# Generate OpenAPI 3.1 docs
echo "Generating API documentation..."
(
	cd "$script_dir/.."
	swag init \
		--v3.1 \
		--generalInfo ./cmd/server/main.go \
		--output "$outPath" \
		--parseDependency
)

echo "API documentation generated in $outPath/"
