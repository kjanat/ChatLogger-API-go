// Package api implements the REST API router and routes for the ChatLogger API.
// This file sets up the Gin router with middleware, route groups, and injects
// dependencies for handlers.
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/middleware"
)

// NewRouter sets up the Gin router with defined routes.
func NewRouter(services *AppServices, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// Apply global middlewares
	router.Use(middleware.VersionHeader())

	// Set up API documentation routes
	setupSwaggerRoutes(router)

	// Add API routes
	addRoutes(router, services, jwtSecret)

	return router
}
