package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FileRoutes(api fiber.Router, db *gorm.DB) {
	fileDownloadRepo := repositories.NewFileDownloadRepository(db)
	fileDownloadController := controllers.NewFileDownloadController(fileDownloadRepo)

	files := api.Group("/files/d", middleware.Auth(db))
	files.Get("/:uid", fileDownloadController.GetFiles)
}
