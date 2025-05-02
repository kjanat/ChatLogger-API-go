// Package api implements the REST API router and routes for the ChatLogger API.
// This file sets up the Gin router with middleware, route groups, and injects
// dependencies for handlers.
package api

import (
	"github.com/gin-gonic/gin"
    "log" //new
    "net/http" //new
    "os" //new

	"github.com/kjanat/chatlogger-api-go/internal/middleware"

    // Import the generated Swagger docs
	"github.com/kjanat/chatlogger-api-go/internal/version"

    docs "github.com/kjanat/chatlogger-api-go/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter sets up the Gin router with defined routes.
func NewRouter(services *AppServices, jwtSecret string) *gin.Engine {
    // Set Swagger info dynamically from version package
    docs.SwaggerInfo.Version = version.Version
    docs.SwaggerInfo.Host = "localhost:8080"
    docs.SwaggerInfo.BasePath = "/api/v1"
    docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()

	// Apply global middlewares
	router.Use(middleware.VersionHeader())
    
    // Add debug handler for swagger docs to diagnose the issue
    router.GET("/docs/swagger.json", func(c *gin.Context) {
        // Try to read the swagger.json file directly
        data, err := os.ReadFile("./docs/swagger.json")
        if err != nil {
            log.Printf("Error reading swagger file: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
                "docs_exists": false,
            })
            return
        }
        
        // If found, return the raw JSON for debugging
        c.Header("Content-Type", "application/json")
        c.String(http.StatusOK, string(data))
    })
    
    // Fix the Swagger documentation route
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
        ginSwagger.URL("/docs/swagger.json"),
        ginSwagger.DefaultModelsExpandDepth(-1),
    ))

	// Add API routes
	addRoutes(router, services, jwtSecret)

	return router
}
