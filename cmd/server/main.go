// Package main is the entry point for the ChatLogger API server.
// It initializes the configuration, sets up database connections,
// initializes all services and repositories, and starts the HTTP server.
package main

import (
	"log"

	"github.com/kjanat/ChatLogger-API-go/internal/api"
	"github.com/kjanat/ChatLogger-API-go/internal/config"
	"github.com/kjanat/ChatLogger-API-go/internal/jobs"
	"github.com/kjanat/ChatLogger-API-go/internal/repository"
	"github.com/kjanat/ChatLogger-API-go/internal/service"
	"github.com/kjanat/ChatLogger-API-go/internal/version"
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

	// 2. Initialize Database with migrations enabled (server is responsible for migrations)
	dbOptions := repository.DefaultDatabaseOptions()
	dbOptions.RunMigrations = true // Server should run migrations

	db, err := repository.NewDatabaseWithOptions(cfg.DatabaseURL, dbOptions)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// 3. Initialize Repositories
	orgRepo := repository.NewOrganizationRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	exportRepo := repository.NewExportRepository(db.DB)

	// 4. Initialize Job Queue
	queue := jobs.NewQueue(cfg.RedisAddr)
	defer func() {
		if err := queue.Close(); err != nil {
			log.Printf("Error closing queue connection: %v", err)
		}
	}()

	// 5. Initialize Services
	orgService := service.NewOrganizationService(orgRepo)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	chatService := service.NewChatService(chatRepo)
	messageService := service.NewMessageService(messageRepo)
	exportService := service.NewExportService(exportRepo, queue)

	// 6. Bundle services for dependency injection
	services := &api.AppServices{
		OrganizationService: orgService,
		APIKeyService:       apiKeyService,
		UserService:         userService,
		ChatService:         chatService,
		MessageService:      messageService,
		ExportService:       exportService,
		Config: &api.AppConfig{
			ExportDir: cfg.ExportDir,
		},
	}

	// 7. Set up Gin Router with routes and inject services
	router := api.NewRouter(services, cfg.JWTSecret)

	// 8. Start the Server
	port := cfg.ServerPort
	log.Printf("Server listening on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
