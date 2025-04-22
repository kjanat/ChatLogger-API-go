$flags = @"
-X ChatLogger-API-go/internal/version.Version=$(bash -c 'grep -oP ''Version = "\K[^"]+'' internal/version/version.go')
-X ChatLogger-API-go/internal/version.BuildTime=$(bash -c 'date -u +"%Y-%m-%dT%H:%M:%SZ"')
-X ChatLogger-API-go/internal/version.GitCommit=$(bash -c 'git rev-parse HEAD')
"@

go build `
    -ldflags $flags `
    -o chatlogger-api.exe `
    ./cmd/server
