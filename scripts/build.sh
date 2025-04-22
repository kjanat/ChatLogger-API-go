# /bin/bash

go build \
    -ldflags \
    -X 'ChatLogger-API-go/internal/version.Version=$(grep -oP 'Version = "\K[^"]+' internal/version/version.go)' \
    -X 'ChatLogger-API-go/internal/version.BuildTime=$(date -u +"%Y-%m-%dT%H:%M:%SZ")' \
    -X 'ChatLogger-API-go/internal/version.GitCommit=$(git rev-parse HEAD)' \
    -o chatlogger-api \
    ./cmd/server
