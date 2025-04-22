package repository

import (
	"fmt"
	"log"
	"time"

	"ChatLogger-API-go/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents a database connection.
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection.
func NewDatabase(dsn string) (*Database, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate the database schema
	if err := db.AutoMigrate(
		&domain.Organization{},
		&domain.APIKey{},
		&domain.User{},
		&domain.Chat{},
		&domain.Message{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")

	return &Database{DB: db}, nil
}

// Close closes the database connection.
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
