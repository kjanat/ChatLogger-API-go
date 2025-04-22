package domain

import (
	"time"
)

// Role represents user access level.
type Role string

// User roles define the permission levels in the system.
const (
	// RoleSuperAdmin has full system access across all organizations.
	RoleSuperAdmin Role = "superadmin" // Full system access
	// RoleAdmin has full access to their organization's data.
	RoleAdmin Role = "admin" // Full organization access
	// RoleUser has access to their own chats and messages.
	RoleUser Role = "user" // Access to own chats/messages
	// RoleViewer has read-only access to their organization's data.
	RoleViewer Role = "viewer" // Read-only access
)

// User represents a system user.
type User struct {
	ID             uint         `gorm:"primaryKey"                    json:"id"`
	Email          string       `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash   string       `gorm:"size:255;not null"             json:"-"`
	Role           Role         `gorm:"size:20;not null"              json:"role"`
	OrganizationID uint         `gorm:"not null;index"                json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"     json:"-"`
	FirstName      string       `gorm:"size:100"                      json:"first_name,omitempty"`
	LastName       string       `gorm:"size:100"                      json:"last_name,omitempty"`
	CreatedAt      time.Time    `                                     json:"created_at"`
	UpdatedAt      time.Time    `                                     json:"updated_at"`
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByOrganizationID(orgID uint, limit, offset int) ([]User, error)
	Update(user *User) error
	Delete(id uint) error
}

// UserService defines the interface for user business logic.
type UserService interface {
	Register(user *User, password string) error
	Login(email, password string) (string, error) // Returns JWT token
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByOrganizationID(orgID uint, limit, offset int) ([]User, error)
	UpdateUser(user *User) error
	ChangePassword(userID uint, currentPassword, newPassword string) error
	DeleteUser(id uint) error
}
