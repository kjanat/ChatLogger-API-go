package domain

// SwaggerService defines the interface for generating and configuring Swagger documentation.
type SwaggerService interface {
	// SetSwaggerInfo updates the OpenAPI/Swagger documentation information.
	// Parameters:
	// - version: API version
	// - scheme: server scheme (http or https)
	// - host: server hostname
	// - port: server port
	SetSwaggerInfo(version, scheme, host, port string)
}
