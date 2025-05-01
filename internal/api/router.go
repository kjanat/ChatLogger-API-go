// Package api implements the REST API router and routes for the ChatLogger API.
// This file sets up the Gin router with middleware, route groups, and injects
// dependencies for handlers.
package api

import (
	"github.com/gin-gonic/gin"

	"github.com/kjanat/chatlogger-api-go/internal/middleware"

	// Import the generated Swagger docs
	_ "github.com/kjanat/chatlogger-api-go/cmd/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // Make sure to use v2 for OpenAPI 3.x support
)

// NewRouter sets up the Gin router with defined routes.
func NewRouter(services *AppServices, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// Add Swagger UI endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply global middlewares
	router.Use(middleware.VersionHeader())

	// Add API routes
	addRoutes(router, services, jwtSecret)

	return router
}
