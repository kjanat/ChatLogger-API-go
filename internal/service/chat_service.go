package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// ChatService implements the domain.ChatService interface.
type ChatService struct {
	chatRepo domain.ChatRepository
}

// NewChatService creates a new chat service.
func NewChatService(chatRepo domain.ChatRepository) domain.ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

// CreateChat creates a new chat.
func (s *ChatService) CreateChat(chat *domain.Chat) error {
	// Set timestamps
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = time.Now()

	// Create the chat
	return s.chatRepo.Create(chat)
}

// GetByID gets a chat by ID.
func (s *ChatService) GetByID(id uint64) (*domain.Chat, error) {
	return s.chatRepo.FindByID(id)
}

// GetByOrganizationID gets chats by organization ID with pagination.
func (s *ChatService) GetByOrganizationID(orgID uint64, limit, offset int) ([]domain.Chat, error) {
	return s.chatRepo.FindByOrganizationID(orgID, limit, offset)
}

// GetByUserID gets chats by user ID with pagination.
func (s *ChatService) GetByUserID(userID uint64, limit, offset int) ([]domain.Chat, error) {
	return s.chatRepo.FindByUserID(userID, limit, offset)
}

// UpdateChat updates a chat.
func (s *ChatService) UpdateChat(chat *domain.Chat) error {
	// Get the existing chat
	existingChat, err := s.chatRepo.FindByID(chat.ID)
	if err != nil {
		return fmt.Errorf("error finding chat: %w", err)
	}

	if existingChat == nil {
		return errors.New("chat not found")
	}

	// Update timestamp
	chat.UpdatedAt = time.Now()

	// Update the chat
	return s.chatRepo.Update(chat)
}

// DeleteChat deletes a chat.
func (s *ChatService) DeleteChat(id uint64) error {
	return s.chatRepo.Delete(id)
}

// GetChatStats gets chat statistics for an organization.
func (s *ChatService) GetChatStats(
	orgID uint64,
	start, end time.Time,
) (map[string]interface{}, error) {
	// Get chat count in date range
	chatCount, err := s.chatRepo.CountByOrgIDAndDateRange(orgID, start, end)
	if err != nil {
		return nil, fmt.Errorf("error getting chat count: %w", err)
	}

	// Get tag statistics
	tagStats, err := s.chatRepo.GetTagStats(orgID)
	if err != nil {
		return nil, fmt.Errorf("error getting tag stats: %w", err)
	}

	// Combine statistics
	stats := map[string]interface{}{
		"total_chats": chatCount,
		"tag_stats":   tagStats,
		"date_range": map[string]string{
			"start": start.Format(time.RFC3339),
			"end":   end.Format(time.RFC3339),
		},
	}

	return stats, nil
}
