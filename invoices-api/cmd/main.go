// cmd/main.go
package main

import (
	"context"
	"invoices-api/config"
	"invoices-api/internal/app"
	"invoices-api/pkg/database"
	"log"
	"os"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database configuration
	dbConfig := &database.DatabaseConfig{
		Host:     cfg.DBHost,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		Port:     cfg.DBPort,
	}

	// Connect to database with retry
	db := database.ConnectDBWithRetry(dbConfig, 5)

	// Create application
	application, err := app.New(db)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Start application
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := application.Start(port); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Wait for shutdown signal
	<-application.WaitForShutdown()

	// Graceful shutdown
	if err := application.Shutdown(context.Background()); err != nil {
		log.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}
}
