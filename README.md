# ChatLogger API (Go)
<!-- ![ChatLogger API](https://raw.githubusercontent.com/yourusername/ChatLogger-API-go/main/assets/logo.png) -->

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/ChatLogger-API-go)](https://goreportcard.com/report/github.com/kjanat/ChatLogger-API-go)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A multi-tenant backend API for logging and managing chat sessions, supporting both authenticated and unauthenticated usage with analytics capabilities.

## 🚀 Features

- **Multi-tenancy**: Separate organizations with isolated data
- **Dual Authentication**: API key for chat plugins and JWT for dashboard users
- **Role-based Access Control**: Superadmin, admin, user, and viewer roles
- **Analytics**: Usage metrics per organization
- **Export Capabilities**: Export chat data in multiple formats (CSV, JSON)
- **Clean Architecture**: Separation of concerns with layered design

## 🛠️ Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/) with PostgreSQL
- **Authentication**: JWT using [golang-jwt/jwt](https://github.com/golang-jwt/jwt/tree/v5)
- **Password Hashing**: bcrypt
- **Configuration**: Environment-based loading
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 14+
- Docker & Docker Compose (optional for containerized deployment)

## 🚀 Quick Start

### Running with Docker

The easiest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/yourusername/ChatLogger-API-go.git
cd ChatLogger-API-go

# Start the application and database
docker-compose up -d
```

The API will be available at `http://localhost:8080`.

### Running Locally

```bash
# Clone the repository
git clone https://github.com/yourusername/ChatLogger-API-go.git
cd ChatLogger-API-go

# Install dependencies
go mod download

# Setup environment variables (copy example and modify)
cp .env.example .env
# Edit .env with your database credentials and settings

# Run migrations
psql -U postgres -d chatlogger -f migrations/001_initial_schema.sql

# Run the server
go run cmd/server/main.go
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
/cmd/server                → Application entry point
/cmd/tools                 → Utility tools
/internal/
  api/                     → Gin router and setup
  config/                  → Configuration loading
  domain/                  → Core models and interfaces
  handler/                 → Request handlers
  hash/                    → Password hashing utilities
  middleware/              → Auth, RBAC, logging middleware
  repository/              → Database access layer
  service/                 → Business logic layer
  strategy/                → Exporters, auth validators
  version/                 → Version information
/migrations/               → SQL schema migrations
/scripts/                  → Utility scripts
```

## 📊 API Endpoints

### Public API (Chat Plugin)

| Method | Endpoint                                    | Description               |
| ------ | ------------------------------------------- | ------------------------- |
| POST   | `/api/v1/orgs/:slug/chats`                  | Create a new chat session |
| POST   | `/api/v1/orgs/:slug/chats/:chatID/messages` | Add a message to a chat   |

### Dashboard API (Authenticated)

| Method | Endpoint                     | Description                |
| ------ | ---------------------------- | -------------------------- |
| POST   | `/auth/login`                | Login and get JWT cookie   |
| GET    | `/users/me`                  | Get current user info      |
| PATCH  | `/users/me`                  | Update current user        |
| GET    | `/orgs/me/apikeys`           | List organization API keys |
| POST   | `/orgs/me/apikeys`           | Create new API key         |
| DELETE | `/orgs/me/apikeys/:id`       | Revoke an API key          |
| GET    | `/analytics/orgs/me/summary` | Get organization analytics |
| POST   | `/exports/orgs/me`           | Trigger data export        |
| GET    | `/exports/orgs/me/:id`       | Download export file       |

## 🔧 Configuration

The application is configured using environment variables. See `.env.example` for available options.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🚢 Deployment

### Version Information

The application includes version information that can be injected during build:

```bash
# Build with version information
docker build -t chatlogger-api \
  --build-arg VERSION=0.1.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  .
```

### CI/CD

The project includes GitHub Actions workflows for CI/CD in `.github/workflows/`.

## 📝 Next Steps and Enhancements

Here are some potential enhancements planned for future development:

- [ ] **Export Features**: Implement the export functionality as a strategy pattern (JSON/CSV)
- [ ] **Pagination**: Add more robust pagination to list endpoints
- [ ] **Testing**: Add unit and integration tests
- [ ] **Documentation**: Generate API documentation with Swagger
- [ ] **Monitoring**: Add logging and metrics collection

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
