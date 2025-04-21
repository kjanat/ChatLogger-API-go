package domain

import (
	"time"
)

// Chat represents a conversation session
type Chat struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	OrganizationID uint         `gorm:"not null;index" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"-"`
	UserID         *uint        `json:"user_id,omitempty"` // Nullable for anonymous chats
	User           *User        `gorm:"foreignKey:UserID" json:"-"`
	Title          string       `gorm:"size:255" json:"title"`
	Tags           string       `gorm:"type:jsonb" json:"tags"`     // JSON array of tags
	Metadata       string       `gorm:"type:jsonb" json:"metadata"` // Additional JSON metadata
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Messages       []Message    `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
}

// ChatRepository defines the interface for chat data operations
type ChatRepository interface {
	Create(chat *Chat) error
	FindByID(id uint) (*Chat, error)
	FindByOrganizationID(orgID uint, limit, offset int) ([]Chat, error)
	FindByUserID(userID uint, limit, offset int) ([]Chat, error)
	Update(chat *Chat) error
	Delete(id uint) error
	// Advanced queries for analytics
	CountByOrgIDAndDateRange(orgID uint, start, end time.Time) (int64, error)
	GetTagStats(orgID uint) (map[string]int64, error)
}

// ChatService defines the interface for chat business logic
type ChatService interface {
	CreateChat(chat *Chat) error
	GetByID(id uint) (*Chat, error)
	GetByOrganizationID(orgID uint, limit, offset int) ([]Chat, error)
	GetByUserID(userID uint, limit, offset int) ([]Chat, error)
	UpdateChat(chat *Chat) error
	DeleteChat(id uint) error
	// Analytics methods
	GetChatStats(orgID uint, start, end time.Time) (map[string]interface{}, error)
}
