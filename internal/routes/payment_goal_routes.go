package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PaymentGoalRoutes(api fiber.Router, db *gorm.DB) {
	paymentGoalRepo := repositories.NewPaymentGoalRepository(db)
	paymentGoalController := controllers.NewPaymentGoalController(paymentGoalRepo)

	paymentGoals := api.Group("/payment-goals", middleware.Auth(db))
	paymentGoals.Get("/", paymentGoalController.Index)
	paymentGoals.Get("/overview", paymentGoalController.Overview)
	paymentGoals.Get("/:id", paymentGoalController.Show)
}
