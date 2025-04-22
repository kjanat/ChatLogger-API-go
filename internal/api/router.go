// Package api handles the API routes, middleware configuration, and API initialization
// for the ChatLogger application. It defines both public and authenticated routes.
package api

import (
	"ChatLogger-API-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter sets up the Gin router with defined routes.
func NewRouter(services *AppServices, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// Apply global middlewares
	router.Use(middleware.VersionHeader())

	// Add API routes
	addRoutes(router, services, jwtSecret)

	return router
}
