package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles chat-related requests.
type ChatHandler struct {
	chatService    domain.ChatService
	messageService domain.MessageService
}

// NewChatHandler creates a new chat handler.
func NewChatHandler(
	chatService domain.ChatService,
	messageService domain.MessageService,
) *ChatHandler {
	return &ChatHandler{
		chatService:    chatService,
		messageService: messageService,
	}
}

// CreateChatRequest represents the request to create a new chat.
type CreateChatRequest struct {
	Title    string               `json:"title"`
	Tags     []string             `json:"tags,omitempty"`
	Metadata *domain.ChatMetadata `json:"metadata,omitempty"` // Use the structured type
	UserID   *uint64              `json:"user_id,omitempty"`  // Optional, for anonymous chats
}

// CreateChat handles the request to create a new chat.
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Get organization ID from context (set by middleware)
	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}
	orgID := orgIDAny.(uint64)

	// Get a userID if available from the JWT context (may not exist for API key auth)
	var userID *uint64
	userIDInterface, exists := c.Get("userID")
	if exists {
		uid := userIDInterface.(uint64)
		userID = &uid
	} else {
		// Use the userID from the request if provided (primarily for API key scenarios)
		userID = req.UserID
	}

	// Create chat object
	chat := &domain.Chat{
		OrganizationID: orgID,
		UserID:         userID,
		Title:          req.Title,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set tags using the helper method
	if err := chat.SetTags(req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tags: " + err.Error()})
		return
	}

	// Set metadata using the helper method
	if err := chat.SetMetadata(req.Metadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process metadata: " + err.Error()})
		return
	}

	// Create the chat
	if err := h.chatService.CreateChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Chat created successfully",
		"chat_id": chat.ID,
	})
}

// GetChatResponse enhances the Chat domain model for API responses.
type GetChatResponse struct {
	*domain.Chat
	ParsedMetadata *domain.ChatMetadata `json:"metadata,omitempty"` // Parsed metadata
	ParsedTags     []string             `json:"tags,omitempty"`     // Parsed tags
}

// GetChat handles the request to get a chat by ID.
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

	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat: " + err.Error()})
		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}
	orgID := orgIDAny.(uint64)

	if chat.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this chat"})
		return
	}

	// Prepare response object
	response := &GetChatResponse{
		Chat: chat,
	}

	// Parse metadata and tags for the response
	metadata, err := chat.GetMetadata()
	if err == nil { // Ignore parsing errors, return raw if needed
		response.ParsedMetadata = metadata
	}
	tags, err := chat.GetTags()
	if err == nil { // Ignore parsing errors
		response.ParsedTags = tags
	}

	// Optionally include messages
	includeMessages := c.Query("include_messages") == "true"
	if includeMessages {
		messages, err := h.messageService.GetByChatID(chat.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages: " + err.Error()})
			return
		}
		chat.Messages = messages
	}

	// Nullify the raw JSON fields in the response to avoid redundancy
	response.Metadata = ""
	response.Tags = ""

	c.JSON(http.StatusOK, response)
}

// ListChats handles the request to list chats for the current organization.
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
	chats, err := h.chatService.GetByOrganizationID(orgID.(uint64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// UpdateChatRequest represents the request to update a chat.
type UpdateChatRequest struct {
	Title    string               `json:"title,omitempty"`
	Tags     []string             `json:"tags,omitempty"`
	Metadata *domain.ChatMetadata `json:"metadata,omitempty"` // Use the structured type
}

// UpdateChat handles the request to update a chat.
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat: " + err.Error()})
		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}
	orgID := orgIDAny.(uint64)

	if chat.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this chat"})
		return
	}

	// Update chat fields if provided
	updated := false
	if req.Title != "" {
		chat.Title = req.Title
		updated = true
	}

	if req.Tags != nil { // Check if tags field was present in the request
		if err := chat.SetTags(req.Tags); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tags: " + err.Error()})
			return
		}
		updated = true
	}

	if req.Metadata != nil { // Check if metadata field was present in the request
		if err := chat.SetMetadata(req.Metadata); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process metadata: " + err.Error()})
			return
		}
		updated = true
	}

	// Update the chat only if something changed
	if updated {
		chat.UpdatedAt = time.Now()
		if err := h.chatService.UpdateChat(chat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat: " + err.Error()})
			return
		}
	}

	// Prepare response similar to GetChat
	response := &GetChatResponse{
		Chat: chat,
	}
	metadata, _ := chat.GetMetadata()
	response.ParsedMetadata = metadata
	tags, _ := chat.GetTags()
	response.ParsedTags = tags
	response.Metadata = ""
	response.Tags = ""

	c.JSON(http.StatusOK, response)
}

// DeleteChat handles the request to delete a chat.
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
	if chat.OrganizationID != orgID.(uint64) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to delete this chat"},
		)
		return
	}

	// Delete the chat
	if err := h.chatService.DeleteChat(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
