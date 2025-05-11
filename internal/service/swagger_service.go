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
func (s *SwaggerService) SetSwaggerInfo(version, scheme, host, port string) {
	serverAddr := host
	if scheme != "" {
		serverAddr = fmt.Sprintf("%s://%s", scheme, host)
	}
	// If the port is not empty and not the default HTTP/HTTPS ports, include it in the address
	// Otherwise, the default ports (80 for HTTP and 443 for HTTPS) are implied
	// and do not need to be included in the address.
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
