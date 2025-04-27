// Package strategy implements the strategy design pattern for various
// pluggable components of the ChatLogger API.
package strategy

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// Exporter defines the interface for export strategies.
// Implementations can export data in different formats (JSON, CSV, etc.).
type Exporter interface {
	// Export takes arbitrary data and exports it to a specific format.
	// It returns the exported data as bytes and any error encountered.
	Export(data interface{}) ([]byte, error)
}

// JSONExporter implements the Exporter interface for JSON format.
type JSONExporter struct{}

// Export exports data to JSON format.
func (j *JSONExporter) Export(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// CSVExporter implements the Exporter interface for CSV format.
type CSVExporter struct{}

// Export exports data to CSV format.
// It expects 'data' to be a map with at least a 'chats' key containing an array of chats.
func (c *CSVExporter) Export(data interface{}) ([]byte, error) {
	// Convert data to map
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data must be a map[string]interface{}")
	}

	// Get chats from data
	chatsInterface, ok := dataMap["chats"]
	if !ok {
		return nil, fmt.Errorf("data must contain a 'chats' key")
	}

	chats, ok := chatsInterface.([]domain.Chat)
	if !ok {
		return nil, fmt.Errorf("'chats' must be an array of domain.Chat")
	}

	// Create CSV writer
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Write header
	header := []string{
		"Chat ID", "Organization ID", "User ID", "Title", "Created At",
		"Message ID", "Role", "Content", "Timestamp", "Token Count", "Latency",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write data
	for _, chat := range chats {
		chatID := fmt.Sprintf("%d", chat.ID)
		orgID := fmt.Sprintf("%d", chat.OrganizationID)
		var userID string
		if chat.UserID != nil {
			userID = fmt.Sprintf("%d", *chat.UserID)
		} else {
			userID = "N/A"
		}
		title := chat.Title
		createdAt := chat.CreatedAt.Format(time.RFC3339)

		// If no messages, write a row for just the chat
		if len(chat.Messages) == 0 {
			row := []string{chatID, orgID, userID, title, createdAt, "", "", "", "", "", ""}
			if err := writer.Write(row); err != nil {
				return nil, err
			}
			continue
		}

		// Write a row for each message in the chat
		for _, message := range chat.Messages {
			messageID := fmt.Sprintf("%d", message.ID)
			role := string(message.Role)
			content := message.Content
			timestamp := message.CreatedAt.Format(time.RFC3339)
			tokenCount := fmt.Sprintf("%d", message.TokenCount)
			latency := fmt.Sprintf("%d", message.Latency)

			row := []string{
				chatID, orgID, userID, title, createdAt,
				messageID, role, content, timestamp, tokenCount, latency,
			}
			if err := writer.Write(row); err != nil {
				return nil, err
			}
		}
	}

	writer.Flush()
	return []byte(sb.String()), nil
}
