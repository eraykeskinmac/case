package database

import (
	"fmt"
	"invoices-api/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}
}

func ConnectDB(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.User, config.Password, config.DBName, config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// GORM sorgu loglarını açalım
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Şemayı doğrulayarak migrasyon yap
	err = db.AutoMigrate(&models.Invoice{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")

	if err := SeedData(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	return db, nil
}

func ConnectDBWithRetry(config *DatabaseConfig, maxRetries int) *gorm.DB {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = ConnectDB(config)
		if err == nil {
			log.Println("Successfully connected to database")
			return db
		}

		log.Printf("Failed to connect to database (Attempt %d/%d): %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			retryDelay := time.Duration(i+1) * 5 * time.Second
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	log.Fatalf("Could not connect to database after %d attempts", maxRetries)
	return nil
}
