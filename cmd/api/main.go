package main

import (
	"fmt"
	"golang-api/internal/config"
	"golang-api/internal/middleware"
	"golang-api/internal/models"
	"golang-api/internal/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	config.ConnectDatabase()

	// Auto migrate models - mirip dengan php artisan migrate
	db := config.GetDB()
	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully!")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
