package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UptimeMonitorRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controllers.NewUptimeMonitorController(db)
	monitors := api.Group("/uptime-monitors", middleware.Auth(db))

	monitors.Get("/", ctrl.Index)
	monitors.Post("/run-checks", ctrl.RunChecks)
	monitors.Post("/", ctrl.Store)
	monitors.Get("/:id", ctrl.Show)
	monitors.Put("/:id", ctrl.Update)
	monitors.Delete("/:id", ctrl.Destroy)
}
