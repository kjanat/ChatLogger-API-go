package service

import (
	"fmt"
	"log"

	docs "github.com/kjanat/chatlogger-api-go/docs"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// SwaggerService implements the domain.SwaggerService interface for managing
// Swagger/OpenAPI documentation.
type SwaggerService struct{}

// NewSwaggerService creates a new Swagger service.
func NewSwaggerService() domain.SwaggerService {
	return &SwaggerService{}
}

// SetSwaggerInfo updates the OpenAPI/Swagger documentation information.
func (s *SwaggerService) SetSwaggerInfo(version, host, port string) {
	serverAddr := host
	if port != "" && port != "80" && port != "443" {
		serverAddr = fmt.Sprintf("%s:%s", host, port)
	}

	log.Printf("Setting Swagger host to %s with version %s", serverAddr, version)

	// Update swagger info for both OpenAPI 2.0 (Swagger) and OpenAPI 3.0
	docs.SwaggerInfoOpenAPI.Version = version
	docs.SwaggerInfoOpenAPI.Host = serverAddr
	docs.SwaggerInfoOpenAPI.BasePath = "/v1"
	docs.SwaggerInfoOpenAPI.Schemes = []string{"http", "https"}
}
