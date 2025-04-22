package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"ChatLogger-API-go/internal/domain"
	"github.com/gin-gonic/gin"
)

// MessageHandler handles message-related requests.
type MessageHandler struct {
	messageService domain.MessageService
	chatService    domain.ChatService
}

// NewMessageHandler creates a new message handler.
func NewMessageHandler(
	messageService domain.MessageService,
	chatService domain.ChatService,
) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		chatService:    chatService,
	}
}

// CreateMessageRequest represents the request to create a new message.
type CreateMessageRequest struct {
	Role       domain.MessageRole `binding:"required" json:"role"`
	Content    string             `binding:"required" json:"content"`
	Metadata   json.RawMessage    `                   json:"metadata,omitempty"`
	Latency    int                `                   json:"latency,omitempty"`
	TokenCount int                `                   json:"token_count,omitempty"`
}

// CreateMessage handles the request to create a new message in a chat.
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	// Get chat ID from URL
	chatID := c.Param("chatID")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})

		return
	}

	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})

		return
	}

	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

		return
	}

	// Get the chat to validate ownership
	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat"})

		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})

		return
	}

	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}

	// Check if the chat belongs to the organization
	if chat.OrganizationID != uint64(orgID.(uint)) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to add messages to this chat"},
		)

		return
	}

	// Create message object
	message := &domain.Message{
		ChatID:     chat.ID,
		Role:       req.Role,
		Content:    req.Content,
		Metadata:   string(req.Metadata),
		Latency:    req.Latency,
		TokenCount: req.TokenCount,
		CreatedAt:  time.Now(),
	}

	// Create the message
	if err := h.messageService.CreateMessage(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Message created successfully",
		"message_id": message.ID,
	})
}

// GetMessages handles the request to get all messages for a chat.
func (h *MessageHandler) GetMessages(c *gin.Context) {
	// Get chat ID from URL
	chatID := c.Param("chatID")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})

		return
	}

	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})

		return
	}

	// Get the chat to validate ownership
	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat"})

		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})

		return
	}

	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}

	// Check if the chat belongs to the organization
	if chat.OrganizationID != uint64(orgID.(uint)) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to view messages in this chat"},
		)

		return
	}

	// Get messages for the chat
	messages, err := h.messageService.GetByChatID(chat.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})

		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetMessageStats handles the request to get message statistics for an organization.
func (h *MessageHandler) GetMessageStats(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}

	// Parse date range parameters
	startStr := c.DefaultQuery("start", "")
	endStr := c.DefaultQuery("end", "")

	var start, end time.Time

	var err error

	if startStr == "" {
		// Default to 30 days ago
		start = time.Now().AddDate(0, 0, -30)
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})

			return
		}
	}

	if endStr == "" {
		// Default to now
		end = time.Now()
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})

			return
		}
	}

	// Get message statistics
	stats, err := h.messageService.GetMessageStats(uint64(orgID.(uint)), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message statistics"})

		return
	}

	c.JSON(http.StatusOK, stats)
}
