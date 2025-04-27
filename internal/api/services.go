package api

import (
	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// AppConfig contains application configuration values
type AppConfig struct {
	ExportDir string
}

// AppServices contains all the services used by the application.
type AppServices struct {
	UserService         domain.UserService
	OrganizationService domain.OrganizationService
	APIKeyService       domain.APIKeyService
	ChatService         domain.ChatService
	MessageService      domain.MessageService
	ExportService       domain.ExportService
	Config              *AppConfig
}
