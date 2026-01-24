package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NotificationRoutes(api fiber.Router, db *gorm.DB) {
	notificationController := controllers.NewNotificationController()

	notifications := api.Group("/notifications", middleware.Auth(db))
	notifications.Put("/settings", notificationController.UpdateSettings)
}
