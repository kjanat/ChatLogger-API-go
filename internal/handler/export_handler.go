package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kjanat/ChatLogger-API-go/internal/domain"
	"github.com/kjanat/ChatLogger-API-go/internal/middleware"
	"github.com/kjanat/ChatLogger-API-go/internal/strategy"

	"github.com/gin-gonic/gin"
)

// ExportHandler handles export-related requests.
type ExportHandler struct {
	exportService  domain.ExportService
	chatService    domain.ChatService
	messageService domain.MessageService
	exportDir      string
}

// NewExportHandler creates a new export handler.
func NewExportHandler(
	exportService domain.ExportService,
	chatService domain.ChatService,
	messageService domain.MessageService,
	exportDir string,
) *ExportHandler {
	return &ExportHandler{
		exportService:  exportService,
		chatService:    chatService,
		messageService: messageService,
		exportDir:      exportDir,
	}
}

// ExportRequest represents the request to export data.
type ExportRequest struct {
	Format string `binding:"required,oneof=json csv" json:"format"`
	Type   string `binding:"required,oneof=chats messages all" json:"type"`
}

// CreateExport handles the request to create an asynchronous export
func (h *ExportHandler) CreateExport(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Convert string format and type to domain types
	var format domain.ExportFormat
	switch req.Format {
	case "json":
		format = domain.ExportFormatJSON
	case "csv":
		format = domain.ExportFormatCSV
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export format"})
		return
	}

	var exportType domain.ExportType
	switch req.Type {
	case "chats":
		exportType = domain.ExportTypeChats
	case "messages":
		exportType = domain.ExportTypeMessages
	case "all":
		exportType = domain.ExportTypeAll
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export type"})
		return
	}

	// Create the export job
	export, err := h.exportService.CreateExport(
		orgID.(uint64),
		userID.(uint64),
		format,
		exportType,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create export: " + err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"export_id": export.ID,
		"status":    export.Status,
		"message":   "Export job created successfully and queued for processing",
	})
}

// GetExport handles the request to get export status
func (h *ExportHandler) GetExport(c *gin.Context) {
	exportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid export ID"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	export, err := h.exportService.GetExport(exportID, orgID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Export not found"})
		return
	}

	c.JSON(http.StatusOK, export)
}

// DownloadExport handles the request to download a completed export
func (h *ExportHandler) DownloadExport(c *gin.Context) {
	exportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid export ID"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	export, err := h.exportService.GetExport(exportID, orgID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Export not found"})
		return
	}

	// Check if export is completed
	if export.Status != domain.ExportStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Export is not ready for download",
			"status": export.Status,
		})
		return
	}

	// Check if file exists
	if export.FilePath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Export file path not found"})
		return
	}

	// Check if file exists on disk
	if _, err := os.Stat(export.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Export file not found on disk"})
		return
	}

	// Set appropriate content type
	contentType := "application/json"
	if export.Format == domain.ExportFormatCSV {
		contentType = "text/csv"
	}

	filename := filepath.Base(export.FilePath)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", contentType)

	c.File(export.FilePath)
}

// ListExports handles the request to list exports for an organization
func (h *ExportHandler) ListExports(c *gin.Context) {
	// Get pagination parameters
	limit := 10
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if val, err := strconv.Atoi(limitParam); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if val, err := strconv.Atoi(offsetParam); err == nil && val >= 0 {
			offset = val
		}
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	exports, err := h.exportService.ListExports(orgID.(uint64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exports": exports})
}

// SyncExport is the original synchronous export method, kept for backward compatibility
func (h *ExportHandler) SyncExport(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
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
