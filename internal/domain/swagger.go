package domain

// SwaggerService defines the interface for generating and configuring Swagger documentation.
type SwaggerService interface {
	// SetSwaggerInfo updates the OpenAPI/Swagger documentation information.
	// Parameters:
	// - version: API version
	// - host: server hostname
	// - port: server port
	SetSwaggerInfo(version, host, port string)
}
