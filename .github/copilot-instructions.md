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
| ------------- | ----------------------------------- |
| Language      | Go (1.21+)                          |
| Web Framework | Gin                                 |
| ORM           | GORM + PostgreSQL                   |
| JWT           | `github.com/golang-jwt/jwt/v5/tree/v5` |
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

- `POST /auth/login` → Sets httpOnly cookie
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
/cmd/server                → main.go entry
/internal/
  api/                     → Gin router
  handler/                 → Request validation + response
  service/                 → Business logic
  repository/              → GORM or SQL access
  domain/                  → Core models + interfaces
  middleware/              → Auth, RBAC, logging
  strategy/                → Exporters, auth validators
  config/                  → Env loader, DI setup
/migrations/               → SQL schema
```

---

## ✅ Expectations

- All handler inputs validated
- Role-guarded endpoints using middleware
- All external deps injected via constructors
- Unit tests on services (mocking repo)
- GitHub-ready project structure

---

**Start with:**

1. Wire base Gin app with routes
2. Set up PostgreSQL with GORM models + migrations
3. Implement login + role middleware
4. Add basic chat + message creation via API key
5. Implement analytics + export interface as strategy pattern
