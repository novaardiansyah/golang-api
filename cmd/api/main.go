package main

import (
	"golang-api/internal/config"
	"golang-api/internal/middleware"
	"golang-api/internal/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.LoadEnv()

	config.ConnectDatabase()

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

	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	routes.SetupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
  
	if os.Getenv("APP_ENV") == "production" {
		addr = "127.0.0.1:" + port
	}

	log.Printf("Server starting on %s...\n", addr)
	log.Fatal(app.Listen(addr))
}
