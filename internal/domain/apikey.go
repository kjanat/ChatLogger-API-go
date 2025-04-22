// Package domain defines the core business entities, interfaces and types
// that represent the application's domain model for the ChatLogger API.
// It contains domain entities, service interfaces, and repository interfaces.
package domain

import (
	"time"
)

// APIKey represents authentication credentials for organization API access.
type APIKey struct {
	ID             uint       `gorm:"primaryKey"                    json:"id"`
	OrganizationID uint       `gorm:"not null"                      json:"organization_id"`
	HashedKey      string     `gorm:"size:255;uniqueIndex;not null" json:"-"` // Hashed, never return raw
	Label          string     `gorm:"size:100;not null"             json:"label"`
	CreatedAt      time.Time  `                                     json:"created_at"`
	RevokedAt      *time.Time `                                     json:"revoked_at,omitempty"`
}

// APIKeyRepository defines the interface for API key data operations.
type APIKeyRepository interface {
	Create(key *APIKey) error
	FindByID(id uint) (*APIKey, error)
	FindByHashedKey(hashedKey string) (*APIKey, error)
	ListByOrganizationID(orgID uint) ([]APIKey, error)
	Revoke(id uint) error
	Delete(id uint) error
}

// APIKeyService defines the interface for API key business logic.
type APIKeyService interface {
	GenerateKey(orgID uint, label string) (string, error) // Returns the raw key (only shown once)
	ValidateKey(rawKey string) (*APIKey, error)
	GetByID(id uint) (*APIKey, error)
	ListByOrganizationID(orgID uint) ([]APIKey, error)
	RevokeKey(id uint) error
	DeleteKey(id uint) error
}
