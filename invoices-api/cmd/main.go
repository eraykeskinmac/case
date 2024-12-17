package main

import (
	"context"
	"errors"
	"invoices-api/internal/docs"
	"invoices-api/internal/handlers"
	"invoices-api/internal/repository"
	"invoices-api/pkg/database"
	"invoices-api/pkg/middleware"
	"invoices-api/pkg/validator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/gofiber/swagger"
)

func setupSwagger(app *fiber.App) {
	swaggerSpec := docs.GenerateSwaggerSpec()

	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		return c.JSON(swaggerSpec)
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger.json",
		DeepLinking: true,
	}))
}

func setupRoutes(app *fiber.App, invoiceHandler *handlers.InvoiceHandler, healthHandler *handlers.HealthHandler) {
	api := app.Group("/api")

	api.Get("/health", healthHandler.Check)

	v1 := api.Group("/v1")
	invoices := v1.Group("/invoices")

	invoices.Get("/", invoiceHandler.GetInvoices)
	invoices.Get("/:id", invoiceHandler.GetInvoiceByID)
	invoices.Post("/", invoiceHandler.CreateInvoice)
	invoices.Put("/:id", invoiceHandler.UpdateInvoice)
	invoices.Delete("/:id", invoiceHandler.DeleteInvoice)

	api.Get("/test-panic", handlers.PanicTestHandler)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		AppName:      "Invoice API v1.0",
	})

	app.Use(middleware.PrometheusMiddleware())
	app.Get("/metrics", middleware.PrometheusHandler())

	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     `{"time":"${time}","pid":"${pid}","status":${status},"method":"${method}","path":"${path}","latency":"${latency}","error":"${error}"}` + "\n",
		TimeFormat: time.RFC3339,
	}))
	app.Use(middleware.RecoverMiddleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	dbConfig := database.NewDatabaseConfig()

	db, err := database.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	invoiceValidator := validator.NewInvoiceValidator()
	invoiceRepo := repository.NewInvoiceRepository(db)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceRepo, invoiceValidator)
	healthHandler := handlers.NewHealthHandler(db)

	setupRoutes(app, invoiceHandler, healthHandler)

	setupSwagger(app)
	log.Println("Swagger documentation is available at /swagger")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server started on port 3000")
		if err := app.Listen(":3000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()

	sig := <-quit
	log.Printf("Caught signal %v, starting graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err == nil {
		log.Println("Closing database connection...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed successfully")
	}

	log.Println("Bye!")
}
