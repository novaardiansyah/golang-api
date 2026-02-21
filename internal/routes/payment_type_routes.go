package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PaymentTypeRoutes(api fiber.Router, db *gorm.DB) {
	paymentTypeRepo := repositories.NewPaymentTypeRepository(db)
	paymentTypeController := controllers.NewPaymentTypeController(paymentTypeRepo)

	paymentTypes := api.Group("/payment-types", middleware.Auth(db))
	paymentTypes.Get("/", paymentTypeController.Index)
}
