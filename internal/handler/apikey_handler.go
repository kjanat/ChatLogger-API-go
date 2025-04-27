// Package handler implements HTTP request handlers for the ChatLogger API.
// It contains handler functions for routing, request validation,
// and response formatting by interacting with service layer components.
package handler

import (
	"fmt"
	"net/http"

	"github.com/kjanat/ChatLogger-API-go/internal/domain"

	"github.com/gin-gonic/gin"
)

// APIKeyHandler handles API key-related requests.
type APIKeyHandler struct {
	apiKeyService domain.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler.
func NewAPIKeyHandler(apiKeyService domain.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// GenerateKeyRequest represents the request to generate a new API key.
type GenerateKeyRequest struct {
	Label string `binding:"required" json:"label"`
}

// GenerateKey handles the request to generate a new API key.
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
	rawKey, err := h.apiKeyService.GenerateKey(orgID.(uint64), req.Label)
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

// ListKeys handles the request to list API keys for an organization.
func (h *APIKeyHandler) ListKeys(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Get API keys
	keys, err := h.apiKeyService.ListByOrganizationID(orgID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list API keys"})
		return
	}

	c.JSON(http.StatusOK, keys)
}

// validateKeyAccess is a helper function to validate API key access and permissions.
// It returns the key ID, API key, and a boolean indicating if validation was successful.
// If validation fails, it sets the appropriate HTTP response and returns false.
func (h *APIKeyHandler) validateKeyAccess(c *gin.Context, actionName string) (uint64, *domain.APIKey, bool) {
	// Get key ID from URL
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key ID is required"})
		return 0, nil, false
	}

	var id uint64
	if _, err := fmt.Sscanf(keyID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return 0, nil, false
	}

	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return 0, nil, false
	}

	// Get the API key
	key, err := h.apiKeyService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API key"})
		return 0, nil, false
	}

	if key == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return 0, nil, false
	}

	// Check if the key belongs to the organization
	if key.OrganizationID != orgID.(uint64) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to " + actionName + " this API key"},
		)
		return 0, nil, false
	}

	return id, key, true
}

// RevokeKey handles the request to revoke an API key.
func (h *APIKeyHandler) RevokeKey(c *gin.Context) {
	id, _, ok := h.validateKeyAccess(c, "revoke")
	if !ok {
		return
	}

	// Revoke the key
	if err := h.apiKeyService.RevokeKey(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked successfully"})
}

// DeleteKey handles the request to delete an API key.
func (h *APIKeyHandler) DeleteKey(c *gin.Context) {
	id, _, ok := h.validateKeyAccess(c, "delete")
	if !ok {
		return
	}

	// Delete the key
	if err := h.apiKeyService.DeleteKey(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}
