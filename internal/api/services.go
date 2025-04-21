package api

import (
	"ChatLogger-API-go/internal/domain"
)

// AppServices contains all the services used by the application
type AppServices struct {
	UserService         domain.UserService
	OrganizationService domain.OrganizationService
	APIKeyService       domain.APIKeyService
	ChatService         domain.ChatService
	MessageService      domain.MessageService
}
