package routes

import (
	"golang-api/internal/config"
	"golang-api/internal/controllers"
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes - mirip dengan routes/api.php di Laravel
func SetupRoutes(app *fiber.App) {
	// Get database instance
	db := config.GetDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)

	// Initialize controllers
	userController := controllers.NewUserController(userRepo)
	productController := controllers.NewProductController(productRepo)

	// API routes group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is running",
		})
	})

	// User routes - RESTful API
	users := api.Group("/users")
	users.Get("/", userController.Index)         // GET /api/users
	users.Get("/:id", userController.Show)       // GET /api/users/:id
	users.Post("/", userController.Store)        // POST /api/users
	users.Put("/:id", userController.Update)     // PUT /api/users/:id
	users.Delete("/:id", userController.Destroy) // DELETE /api/users/:id

	// Product routes - RESTful API
	products := api.Group("/products")
	products.Get("/", productController.Index)         // GET /api/products
	products.Get("/:id", productController.Show)       // GET /api/products/:id
	products.Post("/", productController.Store)        // POST /api/products
	products.Put("/:id", productController.Update)     // PUT /api/products/:id
	products.Delete("/:id", productController.Destroy) // DELETE /api/products/:id
}
