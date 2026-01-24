package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PaymentAccountRoutes(api fiber.Router, db *gorm.DB) {
	paymentAccountRepo := repositories.NewPaymentAccountRepository(db)
	paymentAccountController := controllers.NewPaymentAccountController(paymentAccountRepo)

	paymentAccounts := api.Group("/payment-accounts", middleware.Auth(db))
	paymentAccounts.Get("/", paymentAccountController.Index)
}
