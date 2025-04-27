package repository

import (
	"errors"
	"time"

	"github.com/kjanat/ChatLogger-API-go/internal/domain"
	"gorm.io/gorm"
)

// MessageRepo implements the domain.MessageRepository interface.
type MessageRepo struct {
	db *Database
}

// NewMessageRepository creates a new message repository.
func NewMessageRepository(db *Database) domain.MessageRepository {
	return &MessageRepo{db: db}
}

// Create creates a new message.
func (r *MessageRepo) Create(message *domain.Message) error {
	return r.db.Create(message).Error
}

// FindByID finds a message by ID.
func (r *MessageRepo) FindByID(id uint64) (*domain.Message, error) {
	var message domain.Message

	err := r.db.First(&message, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &message, nil
}

// FindByChatID finds messages by chat ID.
func (r *MessageRepo) FindByChatID(chatID uint64) ([]domain.Message, error) {
	var messages []domain.Message
	err := r.db.Where("chat_id = ?", chatID).Order("created_at ASC").Find(&messages).Error

	return messages, err
}

// CountByOrgIDAndDateRange counts messages in a date range for an organization.
func (r *MessageRepo) CountByOrgIDAndDateRange(orgID uint64, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Message{}).
		Joins("JOIN chats ON messages.chat_id = chats.id").
		Where("chats.organization_id = ? AND messages.created_at BETWEEN ? AND ?", orgID, start, end).
		Count(&count).Error

	return count, err
}

// GetRoleStats gets statistics for message roles in an organization.
func (r *MessageRepo) GetRoleStats(orgID uint64) (map[domain.MessageRole]int64, error) {
	type Result struct {
		Role  domain.MessageRole
		Count int64
	}

	var results []Result

	err := r.db.Model(&domain.Message{}).
		Select("messages.role, COUNT(*) as count").
		Joins("JOIN chats ON messages.chat_id = chats.id").
		Where("chats.organization_id = ?", orgID).
		Group("messages.role").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[domain.MessageRole]int64)
	for _, r := range results {
		stats[r.Role] = r.Count
	}

	return stats, nil
}

// GetLatencyStats gets latency statistics for an organization.
func (r *MessageRepo) GetLatencyStats(orgID uint64) (map[string]float64, error) {
	type Result struct {
		Min   float64
		Max   float64
		Avg   float64
		Total int64
	}

	var result Result

	err := r.db.Model(&domain.Message{}).
		Select("MIN(messages.latency) as min, MAX(messages.latency) as max, AVG(messages.latency) as avg, COUNT(*) as total").
		Joins("JOIN chats ON messages.chat_id = chats.id").
		Where("chats.organization_id = ? AND messages.role = ?", orgID, domain.MessageRoleAssistant).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	stats := map[string]float64{
		"min": result.Min,
		"max": result.Max,
		"avg": result.Avg,
	}

	return stats, nil
}

// GetTokenCountStats gets token count statistics for an organization.
func (r *MessageRepo) GetTokenCountStats(orgID uint64) (map[string]int64, error) {
	type Result struct {
		Total int64
		Avg   float64
	}

	var result Result

	err := r.db.Model(&domain.Message{}).
		Select("SUM(messages.token_count) as total, AVG(messages.token_count) as avg").
		Joins("JOIN chats ON messages.chat_id = chats.id").
		Where("chats.organization_id = ?", orgID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	stats := map[string]int64{
		"total": result.Total,
		"avg":   int64(result.Avg),
	}

	return stats, nil
}
