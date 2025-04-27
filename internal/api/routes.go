// Package api implements the REST API router and routes for the ChatLogger API.
// This file defines the route setup for different API endpoints, including public
// routes for chat plugins, authenticated routes for dashboard users, and admin routes.
package api

import (
	"net/http"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/handler"
	"github.com/kjanat/chatlogger-api-go/internal/middleware"
	"github.com/kjanat/chatlogger-api-go/internal/version"

	"github.com/gin-gonic/gin"
)

// addRoutes adds API routes to the router.
func addRoutes(router *gin.Engine, services *AppServices, jwtSecret string) {
	// Public health and version endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":    version.Version,
			"build_time": version.BuildTime,
			"git_commit": version.GitCommit,
		})
	})

	// Auth routes (no auth required)
	authHandler := handler.NewAuthHandler(services.UserService)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register) // In a real app, this might be admin-only
		authGroup.POST("/logout", authHandler.Logout)
	}

	// API routes for the Dashboard (JWT auth required)
	dashboardGroup := router.Group("/api/v1")
	dashboardGroup.Use(middleware.JWTAuth(jwtSecret))
	{
		// User routes
		userHandler := handler.NewUserHandler(services.UserService)
		userGroup := dashboardGroup.Group("/users")
		{
			userGroup.GET("/me", userHandler.GetMe)
			userGroup.PATCH("/me", userHandler.UpdateMe)
			userGroup.POST("/me/password", userHandler.ChangePassword)
		}

		// Organization routes - admin access only
		orgGroup := dashboardGroup.Group("/orgs/me")
		orgGroup.Use(middleware.RoleRequired(domain.RoleAdmin))
		{
			// API key management
			apiKeyHandler := handler.NewAPIKeyHandler(services.APIKeyService)
			orgGroup.GET("/apikeys", apiKeyHandler.ListKeys)
			orgGroup.POST("/apikeys", apiKeyHandler.GenerateKey)
			orgGroup.DELETE("/apikeys/:id", apiKeyHandler.RevokeKey)
		}

		// Chat routes - any authenticated user
		chatHandler := handler.NewChatHandler(services.ChatService, services.MessageService)
		chatGroup := dashboardGroup.Group("/chats")
		{
			chatGroup.POST("", chatHandler.CreateChat)
			chatGroup.GET("", chatHandler.ListChats)
			chatGroup.GET("/:chatID", chatHandler.GetChat)
			chatGroup.PATCH("/:chatID", chatHandler.UpdateChat)
			chatGroup.DELETE("/:chatID", chatHandler.DeleteChat)
		}

		// Message routes - any authenticated user
		messageHandler := handler.NewMessageHandler(services.MessageService, services.ChatService)
		dashboardGroup.GET("/chats/:chatID/messages", messageHandler.GetMessages)

		// Analytics routes
		dashboardGroup.GET("/analytics/messages", messageHandler.GetMessageStats)

		// Export routes - available to all authenticated users
		exportHandler := handler.NewExportHandler(
			services.ExportService,
			services.ChatService,
			services.MessageService,
			services.Config.ExportDir,
		)

		// Async export endpoints
		dashboardGroup.POST("/exports", exportHandler.CreateExport)
		dashboardGroup.GET("/exports", exportHandler.ListExports)
		dashboardGroup.GET("/exports/:id", exportHandler.GetExport)
		dashboardGroup.GET("/exports/:id/download", exportHandler.DownloadExport)

		// Legacy sync export endpoint (for backward compatibility)
		dashboardGroup.POST("/exports/sync", exportHandler.SyncExport)
	}

	// Public API routes (API key auth required)
	publicAPIGroup := router.Group("/api/v1/orgs/:slug")
	publicAPIGroup.Use(middleware.APIKeyAuth(services.APIKeyService))
	publicAPIGroup.Use(middleware.ValidateSlugAccess(services.OrganizationService))
	{
		// Chat and message creation for external integrations
		chatHandler := handler.NewChatHandler(services.ChatService, services.MessageService)
		messageHandler := handler.NewMessageHandler(services.MessageService, services.ChatService)

		publicAPIGroup.POST("/chats", chatHandler.CreateChat)
		publicAPIGroup.POST("/chats/:chatID/messages", messageHandler.CreateMessage)
	}
}
