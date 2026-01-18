package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"mbkm-go/config"
	"mbkm-go/database"
	"mbkm-go/internal/middleware"
	"mbkm-go/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Optional: Run migrations (uncomment if you want auto-migration)
	// if err := database.Migrate(); err != nil {
	// 	log.Printf("Warning: Migration failed: %v", err)
	// }

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      config.AppConfig.AppName,
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(middleware.CORS())

	// Static files for uploads
	app.Static("/uploads", "./uploads")

	// Setup routes
	routes.SetupRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"app":     config.AppConfig.AppName,
			"version": "1.0.0",
		})
	})

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start server
	port := config.AppConfig.AppPort
	log.Printf("🚀 Server starting on http://localhost:%s", port)
	log.Printf("📚 API endpoints available at http://localhost:%s/api/v1", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
