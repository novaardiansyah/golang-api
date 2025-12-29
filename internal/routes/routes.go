package routes

import (
	"golang-api/internal/config"
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	userRepo := repositories.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo)
	authController := controllers.NewAuthController(*userRepo)

	api := app.Group("/api")

	api.Get("/documentation/*", swagger.HandlerDefault)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is running",
		})
	})

	api.Get("/test-email", middleware.Auth(), func(c *fiber.Ctx) error {
		to := c.Query("to", "admin@novadev.my.id")

		err := utils.SendEmail(to, "Testing Golang Email", map[string]any{
			"Name":    "User Testing",
			"Message": "This is a dynamic message from the route!",
		}, "deploy/resources/views/emails/main.html", "deploy/resources/views/emails/test_content.html")

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to send email: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Email has been sent to " + to,
		})
	})

	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)

	users := api.Group("/users", middleware.Auth())
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)

	galleryRepo := repositories.NewGalleryRepository(db)
	galleryController := controllers.NewGalleryController(galleryRepo)

	galleries := api.Group("/galleries", middleware.Auth())
	galleries.Get("/", galleryController.Index)
	galleries.Post("/upload", galleryController.Upload)

	paymentRepo := repositories.NewPaymentRepository(db)
	paymentController := controllers.NewPaymentController(paymentRepo)

	payments := api.Group("/payments", middleware.Auth())
	payments.Get("/", paymentController.Index)
	payments.Get("/summary", paymentController.Summary)
	payments.Get("/:id", paymentController.Show)
	payments.Get("/:id/attachments", paymentController.GetAttachments)

	notificationController := controllers.NewNotificationController()

	notifications := api.Group("/notifications", middleware.Auth())
	notifications.Put("/settings", notificationController.UpdateSettings)
}
