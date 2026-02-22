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
	payments.Post("/generate-report", paymentController.GenerateReport)

	payments.Get("/:id", paymentController.Show)
	payments.Get("/:id/items/summary", paymentController.GetItemsSummary)
	payments.Get("/:id/items/attached", paymentController.GetItemsAttached)
	payments.Get("/:id/items/not-attached", paymentController.GetItemsNotAttached)
}
