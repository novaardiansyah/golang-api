/*
 * Project Name: routes
 * File: uptime_monitor_log_routes.go
 * Created Date: Wednesday February 11th 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UptimeMonitorLogRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controllers.NewUptimeMonitorLogController(db)
	logs := api.Group("/uptime-monitor-logs", middleware.Auth(db))

	logs.Get("/", ctrl.Index)
	logs.Post("/", ctrl.Store)
	logs.Get("/:id", ctrl.Show)
	logs.Delete("/:id", ctrl.Destroy)
}
