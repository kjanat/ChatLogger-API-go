FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application with version information
ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT

# If build args aren't provided, use defaults from the code
RUN if [ -z "$VERSION" ]; then \
        VERSION=$(grep -oP 'Version = "\K[^"]+' internal/version/version.go); \
    fi && \
    if [ -z "$BUILD_TIME" ]; then \
        BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
    fi && \
    if [ -z "$GIT_COMMIT" ]; then \
        GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown"); \
    fi && \
    echo "Building version: $VERSION, build time: $BUILD_TIME, git commit: $GIT_COMMIT" && \
    go build -ldflags "-X 'ChatLogger-API-go/internal/version.Version=$VERSION' -X 'ChatLogger-API-go/internal/version.BuildTime=$BUILD_TIME' -X 'ChatLogger-API-go/internal/version.GitCommit=$GIT_COMMIT'" -o /app/chatlogger-api ./cmd/server

# Use a smaller image for the final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    mkdir -p /app/migrations

# Copy the binary and migrations from the builder stage
COPY --from=builder /app/chatlogger-api /app/
COPY --from=builder /app/migrations /app/migrations

# Expose the port the app runs on
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Run the application
CMD ["/app/chatlogger-api"]