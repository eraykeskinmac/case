package app

import (
	"context"
	"fmt"
	"invoices-api/internal/docs"
	"invoices-api/internal/handlers"
	"invoices-api/internal/repository"
	"invoices-api/pkg/middleware"
	"invoices-api/pkg/validator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/gofiber/swagger"
	"gorm.io/gorm"
)

const (
	shutdownTimeout = 30 * time.Second
)

type App struct {
	fiber    *fiber.App
	db       *gorm.DB
	shutdown chan os.Signal
}

func New(db *gorm.DB) (*App, error) {
	app := &App{
		db:       db,
		shutdown: make(chan os.Signal, 1),
	}

	if err := app.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize application: %w", err)
	}

	return app, nil
}

func (a *App) initialize() error {
	a.fiber = fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		AppName:      "Invoice API v1.0",
	})

	a.setupMiddleware()

	if err := a.setupHandlers(); err != nil {
		return fmt.Errorf("failed to setup handlers: %w", err)
	}

	a.setupSwagger()

	return nil
}

func (a *App) setupMiddleware() {
	a.fiber.Use(middleware.PrometheusMiddleware())
	a.fiber.Get("/metrics", middleware.PrometheusHandler())

	a.fiber.Use(fiberlogger.New(fiberlogger.Config{
		Format:     `{"time":"${time}","pid":"${pid}","status":${status},"method":"${method}","path":"${path}","latency":"${latency}","error":"${error}"}` + "\n",
		TimeFormat: time.RFC3339,
	}))

	a.fiber.Use(middleware.RecoverMiddleware())

	a.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
}

func (a *App) setupHandlers() error {
	validator := validator.NewInvoiceValidator()
	repo := repository.NewInvoiceRepository(a.db)
	invoiceHandler := handlers.NewInvoiceHandler(repo, validator)
	healthHandler := handlers.NewHealthHandler(a.db)

	api := a.fiber.Group("/api")
	api.Get("/health", healthHandler.Check)

	v1 := api.Group("/v1")
	invoices := v1.Group("/invoices")
	{
		invoices.Get("/", invoiceHandler.GetInvoices)
		invoices.Get("/:id", invoiceHandler.GetInvoiceByID)
		invoices.Post("/", invoiceHandler.CreateInvoice)
		invoices.Put("/:id", invoiceHandler.UpdateInvoice)
		invoices.Delete("/:id", invoiceHandler.DeleteInvoice)
	}

	return nil
}

func (a *App) setupSwagger() {
	swaggerSpec := docs.GenerateSwaggerSpec()

	a.fiber.Get("/swagger.json", func(c *fiber.Ctx) error {
		return c.JSON(swaggerSpec)
	})

	a.fiber.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger.json",
		DeepLinking: true,
	}))

	fmt.Println("Swagger documentation is available at /swagger")
}

func (a *App) Start(port string) error {
	signal.Notify(a.shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server started on port %s\n", port)
		if err := a.fiber.Listen(fmt.Sprintf(":%s", port)); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	return nil
}

func (a *App) WaitForShutdown() <-chan os.Signal {
	return a.shutdown
}

func (a *App) Shutdown(ctx context.Context) error {
	fmt.Println("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.fiber.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("error during server shutdown: %w", err)
	}

	fmt.Println("Server shutdown completed successfully")
	return nil
}
