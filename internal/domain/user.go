package domain

import (
	"time"
)

// Role represents user permission levels.
type Role string

const (
	RoleSuperAdmin Role = "superadmin" // Can manage all orgs
	RoleAdmin      Role = "admin"      // Can manage own org
	RoleUser       Role = "user"       // Regular user
	RoleViewer     Role = "viewer"     // Read-only user
)

// User represents a registered user in the system.
type User struct {
	ID             uint64       `gorm:"primaryKey"                json:"id"`
	OrganizationID uint64       `gorm:"not null;index"            json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"-"`
	Email          string       `gorm:"size:255;uniqueIndex"      json:"email"`
	PasswordHash   string       `gorm:"size:255"                  json:"-"`
	Role           Role         `gorm:"size:50;not null"          json:"role"`
	FirstName      string       `gorm:"size:100"                  json:"first_name"`
	LastName       string       `gorm:"size:100"                  json:"last_name"`
	CreatedAt      time.Time    `                                 json:"created_at"`
	UpdatedAt      time.Time    `                                 json:"updated_at"`
	LastLoginAt    *time.Time   `                                 json:"last_login_at,omitempty"`
	Chats          []Chat       `gorm:"foreignKey:UserID"         json:"-"`
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *User) error
	FindByID(id uint64) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByOrganizationID(orgID uint64, limit, offset int) ([]User, error)
	Update(user *User) error
	Delete(id uint64) error
}

// UserService defines the interface for user business logic.
type UserService interface {
	Authenticate(email, password string) (*User, string, error) // Returns user, JWT token, error
	Register(user *User, password string) error
	GetByID(id uint64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByOrganizationID(orgID uint64, limit, offset int) ([]User, error)
	UpdateUser(user *User) error
	ChangePassword(userID uint64, currentPassword, newPassword string) error
	DeleteUser(id uint64) error
}
