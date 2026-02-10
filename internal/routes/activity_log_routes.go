package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"golang-api/internal/controllers"
)

func ActivityLogRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controllers.NewActivityLogController(db)
	
	api.Get("/activity-logs", ctrl.Index)
}
