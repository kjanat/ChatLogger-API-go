package domain

import (
	"time"
)

// Organization represents a tenant in the multi-tenant system.
type Organization struct {
	ID        uint64    `gorm:"primaryKey"                                json:"id"`
	Name      string    `gorm:"size:100;not null"                         json:"name"`
	Slug      string    `gorm:"size:50;uniqueIndex:idx_org_slug;not null" json:"slug"`
	Settings  string    `gorm:"type:jsonb"                                json:"settings"`
	CreatedAt time.Time `                                                 json:"created_at"`
	UpdatedAt time.Time `                                                 json:"updated_at"`
	APIKeys   []APIKey  `gorm:"foreignKey:OrganizationID"                 json:"-"`
	Users     []User    `gorm:"foreignKey:OrganizationID"                 json:"-"`
	Chats     []Chat    `gorm:"foreignKey:OrganizationID"                 json:"-"`
}

// OrganizationRepository defines the interface for organization data operations.
type OrganizationRepository interface {
	Create(org *Organization) error
	FindByID(id uint64) (*Organization, error)
	FindBySlug(slug string) (*Organization, error)
	Update(org *Organization) error
	Delete(id uint64) error
	List(limit, offset int) ([]Organization, error)
}

// OrganizationService defines the interface for organization business logic.
type OrganizationService interface {
	Create(org *Organization) error
	GetByID(id uint64) (*Organization, error)
	GetBySlug(slug string) (*Organization, error)
	Update(org *Organization) error
	Delete(id uint64) error
	List(limit, offset int) ([]Organization, error)
}
