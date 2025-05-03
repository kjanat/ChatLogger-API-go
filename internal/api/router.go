// Package api implements the REST API router and routes for the ChatLogger API.
// This file sets up the Gin router with middleware, route groups, and injects
// dependencies for handlers.
package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/config"
	"github.com/kjanat/chatlogger-api-go/internal/middleware"
	"github.com/kjanat/chatlogger-api-go/internal/version"

	docs "github.com/kjanat/chatlogger-api-go/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter sets up the Gin router with defined routes.
func NewRouter(services *AppServices, jwtSecret string, config *config.Config) *gin.Engine {
	router := gin.Default()

	// Apply global middlewares
	router.Use(middleware.VersionHeader())

	// Set up API documentation routes
	setupSwaggerRoutes(router, version.Version, config.ServerHost, config.ServerPort)

	// Add API routes
	addRoutes(router, services, jwtSecret)

	return router
}

// setupSwaggerRoutes configures the OpenAPI/Swagger documentation routes.
func setupSwaggerRoutes(router *gin.Engine, version, host, port string) {
	log.Printf("Starting ChatLogger API v%s (host: %s, port: %s)",
		version, host, port)

	docs.SwaggerInfoOpenAPI.Version = version
	docs.SwaggerInfoOpenAPI.Host = host + ":" + port
	docs.SwaggerInfoOpenAPI.BasePath = "/api/v1"
	docs.SwaggerInfoOpenAPI.Schemes = []string{"http", "https"}
	docs.SwaggerInfoOpenAPI.InfoInstanceName = "OpenAPI"

	// Static files for documentation
	router.StaticFile("/docs/api.json", "./docs/OpenAPI_swagger.json")
	router.StaticFile("/docs/api.yaml", "./docs/OpenAPI_swagger.yaml")

	router.GET("/openapi/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.DocExpansion("list"),
			ginSwagger.URL("/docs/api.json"),
		),
	)

	router.GET("/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.DocExpansion("list"),
			ginSwagger.URL("/docs/api.json"),
		),
	)
}
