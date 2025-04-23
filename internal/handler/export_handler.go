package handler

import (
	"net/http"
	"strconv"
	"time"

	"ChatLogger-API-go/internal/domain"
	"ChatLogger-API-go/internal/strategy"

	"github.com/gin-gonic/gin"
)

// ExportHandler handles export-related requests.
type ExportHandler struct {
	chatService    domain.ChatService
	messageService domain.MessageService
}

// NewExportHandler creates a new export handler.
func NewExportHandler(
	chatService domain.ChatService,
	messageService domain.MessageService,
) *ExportHandler {
	return &ExportHandler{
		chatService:    chatService,
		messageService: messageService,
	}
}

// ExportRequest represents the request to export data.
type ExportRequest struct {
	Format string `binding:"required,oneof=json csv" json:"format"`
	Type   string `binding:"required,oneof=chats messages all" json:"type"`
}

// Export handles the request to export data.
func (h *ExportHandler) Export(c *gin.Context) {
	var req ExportRequest
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

	// Get chats for the organization
	// Using a reasonable limit for direct export
	chats, err := h.chatService.GetByOrganizationID(orgID.(uint64), 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chats"})
		return
	}

	// For each chat, get its messages if needed
	if req.Type == "all" || req.Type == "messages" {
		for i := range chats {
			messages, err := h.messageService.GetByChatID(chats[i].ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
				return
			}
			chats[i].Messages = messages
		}
	}

	// Prepare data for export
	data := gin.H{
		"organization_id": orgID,
		"export_date":     time.Now().Format(time.RFC3339),
		"export_type":     req.Type,
		"chats":           chats,
	}

	// Select the appropriate exporter based on format
	var exporter strategy.Exporter
	switch req.Format {
	case "json":
		exporter = &strategy.JSONExporter{}
	case "csv":
		exporter = &strategy.CSVExporter{}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export format"})
		return
	}

	// Export the data
	exportData, err := exporter.Export(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export data: " + err.Error()})
		return
	}

	// Set the appropriate content type and filename
	var contentType, extension string
	switch req.Format {
	case "json":
		contentType = "application/json"
		extension = "json"
	case "csv":
		contentType = "text/csv"
		extension = "csv"
	}

	filename := "chatlogger_export_" + strconv.FormatUint(orgID.(uint64), 10) + "_" +
		time.Now().Format("20060102150405") + "." + extension

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, exportData)
}
