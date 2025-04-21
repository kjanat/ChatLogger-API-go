package handler

import (
	"ChatLogger-API-go/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	chatService    domain.ChatService
	messageService domain.MessageService
}

// NewChatHandler creates a new chat handler
func NewChatHandler(chatService domain.ChatService, messageService domain.MessageService) *ChatHandler {
	return &ChatHandler{
		chatService:    chatService,
		messageService: messageService,
	}
}

// CreateChatRequest represents the request to create a new chat
type CreateChatRequest struct {
	Title    string          `json:"title"`
	Tags     []string        `json:"tags,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
	UserID   *uint           `json:"user_id,omitempty"` // Optional, for anonymous chats
}

// CreateChat handles the request to create a new chat
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get organization ID from context (set by middleware)
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Get a userID if available from the JWT context (may not exist for API key auth)
	var userID *uint
	userIDInterface, exists := c.Get("userID")
	if exists {
		uid := userIDInterface.(uint)
		userID = &uid
	} else {
		// Use the userID from the request if provided
		userID = req.UserID
	}

	// Convert tags to JSON
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tags"})
		return
	}

	// Create chat object
	chat := &domain.Chat{
		OrganizationID: orgID.(uint),
		UserID:         userID,
		Title:          req.Title,
		Tags:           string(tagsJSON),
		Metadata:       string(req.Metadata),
	}

	// Create the chat
	if err := h.chatService.CreateChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Chat created successfully",
		"chat_id": chat.ID,
	})
}

// GetChat handles the request to get a chat by ID
func (h *ChatHandler) GetChat(c *gin.Context) {
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

	// Get the chat
	chat, err := h.chatService.GetByID(uint(id))
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
	if chat.OrganizationID != orgID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this chat"})
		return
	}

	// Check if includeMessages query parameter is set
	includeMessages := c.Query("include_messages") == "true"
	if includeMessages {
		// Get messages for the chat
		messages, err := h.messageService.GetByChatID(chat.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
			return
		}
		chat.Messages = messages
	}

	c.JSON(http.StatusOK, chat)
}

// ListChats handles the request to list chats for the current organization
func (h *ChatHandler) ListChats(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get chats
	chats, err := h.chatService.GetByOrganizationID(orgID.(uint), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// UpdateChatRequest represents the request to update a chat
type UpdateChatRequest struct {
	Title    string          `json:"title,omitempty"`
	Tags     []string        `json:"tags,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// UpdateChat handles the request to update a chat
func (h *ChatHandler) UpdateChat(c *gin.Context) {
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

	var req UpdateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get the existing chat
	chat, err := h.chatService.GetByID(uint(id))
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
	if chat.OrganizationID != orgID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this chat"})
		return
	}

	// Update chat fields if provided
	if req.Title != "" {
		chat.Title = req.Title
	}

	if req.Tags != nil {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tags"})
			return
		}
		chat.Tags = string(tagsJSON)
	}

	if len(req.Metadata) > 0 {
		chat.Metadata = string(req.Metadata)
	}

	// Update the chat
	if err := h.chatService.UpdateChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// DeleteChat handles the request to delete a chat
func (h *ChatHandler) DeleteChat(c *gin.Context) {
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

	// Get the chat
	chat, err := h.chatService.GetByID(uint(id))
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
	if chat.OrganizationID != orgID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this chat"})
		return
	}

	// Delete the chat
	if err := h.chatService.DeleteChat(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
