# ChatLogger API (Go)
<!-- ![ChatLogger API](https://raw.githubusercontent.com/kjanat/chatlogger-api-go/master/assets/logo.png) -->

[![Go version](https://img.shields.io/github/go-mod/go-version/kjanat/chatlogger-api-go?logo=Go&logoColor=white)](go.mod)
[![Go Doc](https://godoc.org/github.com/kjanat/chatlogger-api-go?status.svg)][Package documentation]
[![Go Report Card](https://goreportcard.com/badge/github.com/kjanat/chatlogger-api-go)][Go report] <!-- WTF Can't I Capitalize GO?... [![Go Report](https://img.shields.io/badge/Go%20report-A+-brightgreen.svg)][Go report] -->
[![Tag](https://img.shields.io/github/v/tag/kjanat/chatlogger-api-go?sort=semver&label=Tag)](https://github.com/kjanat/chatlogger-api-go/tags)
[![Release Date](https://img.shields.io/github/release-date/kjanat/chatlogger-api-go?label=Release%20date)][Latest release]
[![License: MIT](https://img.shields.io/github/license/kjanat/chatlogger-api-go?label=License)](LICENSE)
[![Commit activity](https://img.shields.io/github/commit-activity/m/kjanat/chatlogger-api-go?label=Commit%20activity)][Commits]
[![Last commit](https://img.shields.io/github/last-commit/kjanat/chatlogger-api-go?label=Last%20commit)][Commits]
[![CI](https://img.shields.io/github/actions/workflow/status/kjanat/chatlogger-api-go/chatlogger-pipeline.yml?logo=github&label=CI)][Build]
[![Codecov](https://img.shields.io/codecov/c/gh/kjanat/chatlogger-api-go?token=QUop5QdCOv&logo=codecov&logoColor=%23F01F7A&label=Coverage)][Codecov] <!-- https://codecov.io/gh/kjanat/chatlogger-api-go/graph/badge.svg?token=QUop5QdCOv -->
[![Issues](https://img.shields.io/github/issues/kjanat/chatlogger-api-go?label=Issues)][Issues]

A multi-tenant backend API for logging and managing chat sessions, supporting both authenticated and unauthenticated usage with analytics capabilities.

## 🚀 Features

- **Multi-tenancy**: Separate organizations with isolated data
- **Dual Authentication**: API key for chat plugins and JWT for dashboard users
- **Role-based Access Control**: Superadmin, admin, user, and viewer roles
- **Analytics**: Usage metrics per organization
- **Export Capabilities**: Export chat data in multiple formats (CSV, JSON) with both sync and async options
- **Clean Architecture**: Separation of concerns with layered design
- **Strategy Pattern**: Used for exporter formats (JSON, CSV) and other pluggable components
- **Asynchronous Jobs**: Background processing with Redis and Asynq

## Documentation

The API is documented using Swagger/OpenAPI:

1. Generate documentation: `./scripts/gendocs.sh`
2. Start the server: `go run cmd/server/main.go`
3. Open browser: http://localhost:8080/swagger/index.html

> [!TIP]
> View the package documentation for more details on the API endpoints and usage.
> [godoc.org][Package documentation]

Note: You must regenerate documentation after making API changes by running the documentation script.

## 🛠️ Tech Stack

- **Language**: [Go 1.24.2+](https://github.com/kjanat/chatlogger-api-go/blob/master/go.mod#L3)
- **Web Framework**: [Gin][Gin]
- **ORM**: [GORM][Gorm] with PostgreSQL
- **Authentication**: JWT using [golang-jwt/jwt][JWT]
- **Password Hashing**: bcrypt
- **Configuration**: Environment-based loading
- **Queue System**: [Asynq][Asynq] with Redis (for async exports)
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 14+
- Redis (optional, for async exports)
- Docker & Docker Compose (optional for containerized deployment)

## 🚀 Quick Start

### Running with Docker

The easiest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/kjanat/chatlogger-api-go.git
cd chatlogger-api-go

# Start the application and database
docker-compose up -d
```

The API will be available at `http://localhost:8080`.

### Running Locally

```bash
# Clone the repository
git clone https://github.com/kjanat/chatlogger-api-go.git
cd chatlogger-api-go

# Install dependencies
go mod download

# Setup environment variables (copy example and modify)
cp .env.example .env
# Edit .env with your database credentials and settings

# Run migrations
psql -U postgres -d chatlogger -f migrations/001_initial_schema.sql
psql -U postgres -d chatlogger -f migrations/002_ensure_defaults.sql
psql -U postgres -d chatlogger -f migrations/003_add_exports_table.sql

# Run the server
go run cmd/server/main.go

# Run the worker (optional, for async exports)
go run cmd/worker/main.go
```

## 🔑 Authentication

### Chat Plugin (Public API)

Uses API key authentication via the `x-organization-api-key` header:

```bash
curl -X POST \
  http://localhost:8080/api/v1/orgs/acme/chats \
  -H 'x-organization-api-key: your-api-key' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Support Chat",
    "tags": ["support", "billing"],
    "metadata": {
      "browser": "Chrome",
      "platform": "Windows"
    }
  }'
```

### Dashboard (User Authentication)

Uses email/password login with JWT authentication:

```bash
# Login and get JWT cookie
curl -X POST \
  http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "admin@example.com",
    "password": "your-password"
  }' \
  -c cookies.txt

# Use the JWT cookie for authenticated requests
curl -X GET \
  http://localhost:8080/users/me \
  -b cookies.txt
```

## 📁 Project Structure

```plaintext
/cmd
  /server                → Main API server entry point
  /tools                 → Utility tools
  /worker                → Background job worker for async exports
/internal
  /api                   → Gin router and setup
  /config                → Configuration loading
  /domain                → Core models and interfaces
  /handler               → Request handlers
  /hash                  → Password hashing utilities
  /jobs                  → Queue and processor for async tasks
  /middleware            → Auth, RBAC, logging middleware
  /repository            → Database access layer
  /service               → Business logic layer
  /strategy              → Strategy pattern implementations (exporters)
  /version               → Version information
/migrations              → SQL schema migrations
/scripts                 → Utility scripts
```

## 📊 API Endpoints

### Public API (Chat Plugin)

| Method   | Endpoint                                     | Description                |
|:---------|:---------------------------------------------|:---------------------------|
| `POST`   | `/api/v1/orgs/:slug/chats`                   | Create a new chat session  |
| `POST`   | `/api/v1/orgs/:slug/chats/:chatID/messages`  | Add a message to a chat    |

### Dashboard API (Authenticated)

| Method   | Endpoint                      | Description                  |
|:---------|:------------------------------|:-----------------------------|
| `POST`   | `/auth/login`                 | Login and get JWT cookie     |
| `GET`    | `/users/me`                   | Get current user info        |
| `PATCH`  | `/users/me`                   | Update current user          |
| `GET`    | `/orgs/me/apikeys`            | List organization API keys   |
| `POST`   | `/orgs/me/apikeys`            | Create new API key           |
| `DELETE` | `/orgs/me/apikeys/:id`        | Revoke an API key            |
| `GET`    | `/analytics/orgs/me/summary`  | Get organization analytics   |

### Export Endpoints

| Method   | Endpoint                      | Description                  |
|:---------|:------------------------------|:-----------------------------|
| `POST`   | `/exports/orgs/me`            | Create async export job      |
| `GET`    | `/exports/orgs/me/:id`        | Get export job status        |
| `GET`    | `/exports/orgs/me/:id/download` | Download completed export  |
| `GET`    | `/exports/orgs/me`            | List export jobs             |
| `POST`   | `/exports/orgs/me/sync`       | Create synchronous export    |

## 🔧 Configuration

The application is configured using environment variables:

| Variable       | Description                      | Default           |
|:---------------|:---------------------------------|:------------------|
| `PORT`         | Server port                      | `8080`            |
| `DATABASE_URL` | PostgreSQL connection string     | Required          |
| `JWT_SECRET`   | Secret for JWT signing           | Required          |
| `REDIS_ADDR`   | Redis address for job queue      | `localhost:6379`  |
| `EXPORT_DIR`   | Directory to store export files  | `./exports`       |

## ⚙️ Export System

The application supports both synchronous and asynchronous exports:

### Synchronous Exports

- Immediate response with download
- Good for small data sets
- Available at `/exports/orgs/me/sync`
- Works without Redis

### Asynchronous Exports

- Queued processing with status tracking
- Better for large data sets
- Requires Redis and worker process
- Creates a job at `/exports/orgs/me`
- Check status at `/exports/orgs/me/:id`
- Download at `/exports/orgs/me/:id/download`

Both export types support JSON and CSV formats.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linting checks
golangci-lint run
```

## 🚢 Deployment

### Version Information

The application includes version information that can be injected during build:

```bash
# Build with version information
docker build -t chatlogger-api \
  --build-arg VERSION=0.3.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  .
```

### Image Verification

All container images are signed using [Sigstore cosign](https://github.com/sigstore/cosign). You can verify the authenticity of the images using the following command:

```bash
# Verify the container image signature
cosign verify \
    --key=cosign.pub \
    ghcr.io/kjanat/chatlogger-api-server:v0.3.0

# Or use the public key URL
cosign verify \
    --key=https://raw.githubusercontent.com/kjanat/chatlogger-api-go/refs/heads/master/cosign.pub \
    ghcr.io/kjanat/chatlogger-api-worker:v0.3.0
```

The public key for verification ([`cosign.pub`](cosign.pub)) is available in the root of the repository. For production deployments, we recommend verifying image signatures as part of your CI/CD pipeline to ensure supply chain security.

### CI/CD

The project includes GitHub Actions workflows for CI/CD:

- **Build Workflow**: Builds, tests, and tags versions on pushes to main
- **Release Workflow**: Creates GitHub releases when tags are pushed

## 📝 Next Steps and Enhancements

Here are some potential enhancements planned for future development:

- [x] **Export Features**: Implemented export functionality as a strategy pattern (JSON/CSV)
- [x] **Async Exports**: Added background processing for large exports
- [ ] **Pagination**: Add more robust pagination to list endpoints
- [ ] **Testing**: Add unit and integration tests
- [ ] **Documentation**: Generate API documentation with Swagger
- [ ] **Monitoring**: Add logging and metrics collection
- [ ] **Real-time notifications**: Add WebSocket support for live updates

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

[Commits]: https://github.com/kjanat/chatlogger-api-go/commits/master/
[Latest release]: https://github.com/kjanat/chatlogger-api-go/releases/latest
[Issues]: https://github.com/kjanat/chatlogger-api-go/issues
[Go report]: https://goreportcard.com/report/github.com/kjanat/chatlogger-api-go
[Gin]: https://github.com/gin-gonic/gin
[Gorm]: https://gorm.io/
[JWT]: https://github.com/golang-jwt/jwt/tree/v5
[Asynq]: https://github.com/hibiken/asynq
[Codecov]: https://codecov.io/gh/kjanat/chatlogger-api-go
[Build]: https://github.com/kjanat/chatlogger-api-go/actions/workflows/chatlogger-pipeline.yml
[Package documentation]: https://godoc.org/github.com/kjanat/chatlogger-api-go
