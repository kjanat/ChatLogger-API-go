package handler

import (
	"net/http"

	"github.com/kjanat/ChatLogger-API-go/internal/domain"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests.
type UserHandler struct {
	userService domain.UserService
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe handles the request to get the current user's info.
func (h *UserHandler) GetMe(c *gin.Context) {
	// Get user ID from context (set by JWTAuth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})

		return
	}

	// Get user from service
	user, err := h.userService.GetByID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateMeRequest represents the update user request body.
type UpdateMeRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateMe handles the request to update the current user's info.
func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})

		return
	}

	// Get the current user
	user, err := h.userService.GetByID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

		return
	}

	// Update the user fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	// Save the updated user
	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePasswordRequest represents the change password request body.
type ChangePasswordRequest struct {
	CurrentPassword string `binding:"required"       json:"current_password"`
	NewPassword     string `binding:"required,min=8" json:"new_password"`
}

// ChangePassword handles the request to change the current user's password.
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})

		return
	}

	// Change the password
	err := h.userService.ChangePassword(userID.(uint64), req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ListOrgUsers handles the request to list all users in the current organization.
func (h *UserHandler) ListOrgUsers(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}

	// Parse pagination parameters
	limit := 20 // Default limit
	offset := 0 // Default offset

	// Get users
	users, err := h.userService.GetByOrganizationID(orgID.(uint64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})

		return
	}

	c.JSON(http.StatusOK, users)
}
