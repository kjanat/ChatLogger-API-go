package main

import (
	"log"

	"ChatLogger-API-go/internal/api"
	"ChatLogger-API-go/internal/config"
	"ChatLogger-API-go/internal/repository"
	"ChatLogger-API-go/internal/service"
	"ChatLogger-API-go/internal/version"
)

func main() {
	// Log version information at startup
	log.Printf("Starting ChatLogger API v%s (built: %s, commit: %s)",
		version.Version, version.BuildTime, version.GitCommit)

	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database
	db, err := repository.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// 3. Initialize Repositories
	orgRepo := repository.NewOrganizationRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// 4. Initialize Services
	orgService := service.NewOrganizationService(orgRepo)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	chatService := service.NewChatService(chatRepo)
	messageService := service.NewMessageService(messageRepo)

	// 5. Bundle services for dependency injection
	services := &api.AppServices{
		OrganizationService: orgService,
		APIKeyService:       apiKeyService,
		UserService:         userService,
		ChatService:         chatService,
		MessageService:      messageService,
	}

	// 6. Set up Gin Router with routes and inject services
	router := api.NewRouter(services, cfg.JWTSecret)

	// 7. Start the Server
	port := cfg.ServerPort
	log.Printf("Server listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
