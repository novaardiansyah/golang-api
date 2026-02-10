package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ActivityLogRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controllers.NewActivityLogController(db)
	activityLogs := api.Group("/activity-logs", middleware.Auth(db))

	activityLogs.Get("/", ctrl.Index)
	activityLogs.Post("/", ctrl.Store)
}
