# 🧱 Project Spec: Chat Logger API Backend

## 🚀 Goal

Build a multi-tenant backend API that:

- Logs and manages chat sessions + messages
- Supports authenticated and unauthenticated usage (chat plugin vs. dashboard)
- Provides usage analytics per organization
- Supports admin/user roles, API key access, and secure login

---

## ⚙️ Stack

| Component     | Choice                              |
|:--------------|:------------------------------------|
| Language      | Go (1.24+)                          |
| Web Framework | Gin                                 |
| ORM           | GORM + PostgreSQL                   |
| JWT           | `github.com/golang-jwt/jwt/tree/v5` |
| Hashing       | bcrypt                              |
| Export Jobs   | Optional: Asynq + Redis             |
| Config        | Viper or env-based loading          |
| Tests         | `stretchr/testify`, mockable layers |
| CI/CD         | GitHub Actions, Docker              |

---

## 📐 Architecture

- Use **Clean Architecture / layered design**
  - `handler` (API layer)
  - `service` (business logic)
  - `repository` (DB access)
  - `domain` (models/interfaces)
- Use **Dependency Injection** (constructor-based)
- Use **Strategy Pattern** for:
  - Export formats (CSV, JSON)
  - API key auth (pluggable validators)

---

## 🧑‍💻 Auth

### A. Chat plugin API (unauthenticated users)

- Auth via `x-organization-api-key` header
- Scopes actions to specific org
- No login required

### B. Dashboard API (admin panel)

- Email/password login
- JWT returned via secure HTTP-only cookie
- Authenticated routes use middleware to extract user + role

### C. Roles

- `superadmin`: Full access
- `admin`: Full org access
- `user`: Own chats/messages
- `viewer`: Optional read-only

---

## 🧱 Core Models

### Organization

- ID, Name, Slug, Settings, APIKeys

### APIKey

- Hashed key, Label, CreatedAt, RevokedAt

### User

- ID, Email, PasswordHash, Role, OrgID

### Chat

- ID, OrgID, UserID (nullable), Title, Tags, Metadata, CreatedAt

### Message

- ChatID, Role (`user`, `assistant`, etc), Content, Metadata, Latency, TokenCount

---

## 📲 Endpoints Overview

### 🔓 Public API (Chat Plugin)

- `POST /api/v1/orgs/:slug/chats`
- `POST /api/v1/orgs/:slug/chats/:chatID/messages`
  - Header: `x-organization-api-key`
  - Body: Content, Role, Metadata

### 🔐 Dashboard API (Users)

- `POST /auth/login` → Sets `httpOnly` cookie
- `GET /users/me`
- `PATCH /users/me`
- Admin:
  - `GET /orgs/me/apikeys`
  - `POST /orgs/me/apikeys`
  - `DELETE /orgs/me/apikeys/:id`

### 📊 Analytics

- `GET /analytics/orgs/me/summary`
  - Messages per day, latency, roles, top users, tags

### 📤 Export

- `POST /exports/orgs/me` → Triggers async/sync export (JSON/CSV)
- `GET /exports/orgs/me/:id` → Download file

---

## 📁 Structure

```plaintext
/cmd
  /server                → main.go entry
  /tools                 → Utility tools (e.g., migrations, data generation)
/internal
  /api                   → Route definitions (Gin)
  /handler               → I/O logic request validation + response
  /service               → Business logic (inject repos, emit events, etc.)
  /repository            → DB access (gorm/sql)
  /domain                → Core models + interfaces
  /middleware            → Auth, RBAC, logging
  /hash                  → Password hashing (bcrypt)
  /version               → Version info (Git commit, build time)
  /strategy              → Plug-and-play logic (exporters, auth validators)
  /config                → Env loader, DI setup
  /jobs                  → Optional async tasks
/migrations              → Starting SQL schema
```

---

## ✅ Expectations

- All handler inputs validated
- Role-guarded endpoints using middleware
- All external deps injected via constructors
- Unit tests on services (mocking repo)
- GitHub-ready project structure

---

Yes — and *hell yes*. You want a **clean architecture with separation of concerns**, dependency injection (DI), and a codebase that’s modular and easy to test or extend. That means **layering** things well and using patterns like:

- **Strategy Pattern** for extensibility (e.g. export formats, API key auth strategies)
- **Dependency Injection** for testability and clear control of resources
- **Interface-driven development** so parts can be mocked, swapped, or extended

Let me give you a full plan that outlines:

1. 📦 **Project layering**  
2. 🔄 **Where to apply the Strategy Pattern**  
3. 🧪 **Where and how to use Dependency Injection**  
4. 🧱 **Interfaces and structure by example**

---

## 🏗️ 1. Project Layering (Clean Architecture)

We'll split the codebase into **layers**, like so.

## 🔄 2. Strategy Pattern — Where To Use It

| Area                       | Strategy Pattern Use Case                        |
| -------------------------- | ------------------------------------------------ |
| 🔑 API Key Auth             | Allow switching between auth strategies per-org  |
| 📤 Exporters                | Different export formats: JSON, CSV, XML, PDF    |
| 📊 Analytics Providers      | You might support internal analytics vs. plug-in |
| 🔁 Rate Limiting Strategies | Per org/user, per endpoint, etc.                 |
| 🌍 Chat Metadata Processors | Extensible pipelines for enrichment, tagging     |

### ✅ Example: Exporter Strategy

```go
type Exporter interface {
    Export(data interface{}) ([]byte, error)
}

type JSONExporter struct{}
type CSVExporter struct{}

func (j *JSONExporter) Export(data interface{}) ([]byte, error) {
    return json.MarshalIndent(data, "", "  ")
}
func (c *CSVExporter) Export(data interface{}) ([]byte, error) {
    // implement csv logic here
}
```

Usage:
```go
var exporter Exporter = getExporterByFormat("csv")
fileBytes, err := exporter.Export(chatData)
```

Let the caller select format dynamically (e.g., via query param or request header).

---

## 🧪 3. Dependency Injection — Where To Apply It

Use DI **explicitly via constructor injection**, not some magic container bullshit.

### Inject in `main.go`

This is where you wire up dependencies:

```go
db := gorm.Open(...)
userRepo := repository.NewUserRepository(db)
chatRepo := repository.NewChatRepository(db)

userService := service.NewUserService(userRepo, ...)
chatService := service.NewChatService(chatRepo, ...)

api := api.NewHandler(userService, chatService, ...)
```

### Key Interfaces to Inject

| Interface        | Used In            | Mock for Tests? |
| ---------------- | ------------------ | --------------- |
| `UserRepository` | `UserService`      | ✅               |
| `ChatRepository` | `ChatService`      | ✅               |
| `Exporter`       | `ExportService`    | ✅               |
| `Mailer`         | `UserService`      | ✅               |
| `KeyValidator`   | `APIKeyMiddleware` | ✅               |
| `TokenManager`   | `AuthMiddleware`   | ✅               |

---

## 🧱 4. Example Interfaces and Structure

### Interface: UserRepository

```go
type UserRepository interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
}
```

### Implementation: `postgres/user_repository.go`

```go
type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

---

### Injecting into Service

```go
type UserService struct {
    repo UserRepository
    mailer Mailer
}

func NewUserService(repo UserRepository, mailer Mailer) *UserService {
    return &UserService{repo, mailer}
}
```

---

### Handler Layer: No Business Logic

```go
func (h *UserHandler) Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    user, token, err := h.service.Authenticate(c.Request.Context(), input.Email, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    setSecureCookie(c, "token", token)
    c.JSON(http.StatusOK, gin.H{"user": user})
}
```

---

## 💣 Extras That Help Testability

| Technique                   | Purpose                            |
| --------------------------- | ---------------------------------- |
| Use `context.Context`       | Pass auth data, request IDs, etc.  |
| Use interfaces for services | Swap/Mock in tests                 |
| Use table-driven tests      | Standard in Go for service/handler |
| Use dependency constructors | Avoid global state/singletons      |
| Avoid direct DB in services | Keep that in repositories only     |

---

## ✅ TL;DR: Contractor Summary

- Use **Gin + GORM + DI via constructors**
- Apply **strategy pattern** for pluggable logic (exporters, auth)
- Structure into **handler → service → repo**
- Use interfaces for **all external concerns**
- Wire dependencies in `main.go` cleanly — don’t hide it behind global context
- Test at **unit level** by mocking interfaces, and **integration level** by spinning up test DB
