package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PaymentRoutes(api fiber.Router, db *gorm.DB) {
	paymentController := controllers.NewPaymentController(db)

	payments := api.Group("/payments", middleware.Auth(db))

	payments.Get("/", paymentController.Index)
	payments.Get("/summary", paymentController.Summary)

	payments.Post("/", paymentController.Store)

	payments.Get("/:id", paymentController.Show)
}
