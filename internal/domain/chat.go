package domain

import (
	"time"
)

// Chat represents a conversation session.
type Chat struct {
	ID             uint64       `gorm:"primaryKey"                json:"id"`
	OrganizationID uint64       `gorm:"not null;index"            json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"-"`
	UserID         *uint64      `                                 json:"user_id,omitempty"` // Nullable for anonymous chats
	User           *User        `gorm:"foreignKey:UserID"         json:"-"`
	Title          string       `gorm:"size:255"                  json:"title"`
	Tags           string       `gorm:"type:jsonb"                json:"tags"`     // JSON array of tags
	Metadata       string       `gorm:"type:jsonb"                json:"metadata"` // Additional JSON metadata
	CreatedAt      time.Time    `                                 json:"created_at"`
	UpdatedAt      time.Time    `                                 json:"updated_at"`
	Messages       []Message    `gorm:"foreignKey:ChatID"         json:"messages,omitempty"`
}

// ChatRepository defines the interface for chat data operations.
type ChatRepository interface {
	Create(chat *Chat) error
	FindByID(id uint64) (*Chat, error)
	FindByOrganizationID(orgID uint64, limit, offset int) ([]Chat, error)
	FindByUserID(userID uint64, limit, offset int) ([]Chat, error)
	Update(chat *Chat) error
	Delete(id uint64) error
	// Advanced queries for analytics
	CountByOrgIDAndDateRange(orgID uint64, start, end time.Time) (int64, error)
	GetTagStats(orgID uint64) (map[string]int64, error)
}

// ChatService defines the interface for chat business logic.
type ChatService interface {
	CreateChat(chat *Chat) error
	GetByID(id uint64) (*Chat, error)
	GetByOrganizationID(orgID uint64, limit, offset int) ([]Chat, error)
	GetByUserID(userID uint64, limit, offset int) ([]Chat, error)
	UpdateChat(chat *Chat) error
	DeleteChat(id uint64) error
	// Analytics methods
	GetChatStats(orgID uint64, start, end time.Time) (map[string]interface{}, error)
}
