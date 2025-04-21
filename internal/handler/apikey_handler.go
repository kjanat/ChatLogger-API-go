package handler

import (
	"ChatLogger-API-go/internal/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyHandler handles API key-related requests
type APIKeyHandler struct {
	apiKeyService domain.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(apiKeyService domain.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// GenerateKeyRequest represents the request to generate a new API key
type GenerateKeyRequest struct {
	Label string `json:"label" binding:"required"`
}

// GenerateKey handles the request to generate a new API key
func (h *APIKeyHandler) GenerateKey(c *gin.Context) {
	var req GenerateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Generate a new API key
	rawKey, err := h.apiKeyService.GenerateKey(orgID.(uint), req.Label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Return the raw key (only shown once)
	c.JSON(http.StatusCreated, gin.H{
		"message": "API key generated successfully",
		"key":     rawKey,
		"label":   req.Label,
	})
}

// ListKeys handles the request to list API keys for an organization
func (h *APIKeyHandler) ListKeys(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Get API keys
	keys, err := h.apiKeyService.ListByOrganizationID(orgID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list API keys"})
		return
	}

	c.JSON(http.StatusOK, keys)
}

// RevokeKey handles the request to revoke an API key
func (h *APIKeyHandler) RevokeKey(c *gin.Context) {
	// Get key ID from URL
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key ID is required"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(keyID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Get the API key
	key, err := h.apiKeyService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API key"})
		return
	}
	if key == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Check if the key belongs to the organization
	if key.OrganizationID != orgID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to revoke this API key"})
		return
	}

	// Revoke the key
	if err := h.apiKeyService.RevokeKey(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked successfully"})
}

// DeleteKey handles the request to delete an API key
func (h *APIKeyHandler) DeleteKey(c *gin.Context) {
	// Get key ID from URL
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key ID is required"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(keyID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Get the API key
	key, err := h.apiKeyService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API key"})
		return
	}
	if key == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Check if the key belongs to the organization
	if key.OrganizationID != orgID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this API key"})
		return
	}

	// Delete the key
	if err := h.apiKeyService.DeleteKey(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}
