package routes

import (
	"golang-api/internal/config"
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	userRepo := repositories.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo)

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is running",
		})
	})

	users := api.Group("/users", middleware.Auth())
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)

	galleryRepo := repositories.NewGalleryRepository(db)
	galleryController := controllers.NewGalleryController(galleryRepo)

	galleries := api.Group("/galleries", middleware.Auth())
	galleries.Get("/", galleryController.Index)

	paymentRepo := repositories.NewPaymentRepository(db)
	paymentController := controllers.NewPaymentController(paymentRepo)

	payments := api.Group("/payments", middleware.Auth())
	payments.Get("/", paymentController.Index)
  payments.Get("/:id", paymentController.Show)
}
