package service

import (
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// MessageService implements the domain.MessageService interface.
type MessageService struct {
	messageRepo domain.MessageRepository
}

// NewMessageService creates a new message service.
func NewMessageService(messageRepo domain.MessageRepository) domain.MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

// CreateMessage creates a new message.
func (s *MessageService) CreateMessage(message *domain.Message) error {
	// Validate the message
	if err := message.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// Set timestamp
	message.CreatedAt = time.Now()

	// Create the message
	return s.messageRepo.Create(message)
}

// GetByID gets a message by ID.
func (s *MessageService) GetByID(id uint64) (*domain.Message, error) {
	return s.messageRepo.FindByID(id)
}

// GetByChatID gets messages by chat ID.
func (s *MessageService) GetByChatID(chatID uint64) ([]domain.Message, error) {
	return s.messageRepo.FindByChatID(chatID)
}

// GetMessageStats gets message statistics for an organization.
func (s *MessageService) GetMessageStats(
	orgID uint64,
	start, end time.Time,
) (map[string]interface{}, error) {
	// Get message count in date range
	messageCount, err := s.messageRepo.CountByOrgIDAndDateRange(orgID, start, end)
	if err != nil {
		return nil, fmt.Errorf("error getting message count: %w", err)
	}

	// Get role statistics
	roleStats, err := s.messageRepo.GetRoleStats(orgID)
	if err != nil {
		return nil, fmt.Errorf("error getting role stats: %w", err)
	}

	// Get latency statistics
	latencyStats, err := s.messageRepo.GetLatencyStats(orgID)
	if err != nil {
		return nil, fmt.Errorf("error getting latency stats: %w", err)
	}

	// Get token count statistics
	tokenStats, err := s.messageRepo.GetTokenCountStats(orgID)
	if err != nil {
		return nil, fmt.Errorf("error getting token stats: %w", err)
	}

	// Combine statistics
	stats := map[string]interface{}{
		"total_messages": messageCount,
		"by_role":        roleStats,
		"latency":        latencyStats,
		"tokens":         tokenStats,
		"date_range": map[string]string{
			"start": start.Format(time.RFC3339),
			"end":   end.Format(time.RFC3339),
		},
	}

	return stats, nil
}
