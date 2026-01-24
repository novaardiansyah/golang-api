package routes

import (
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoutes(api fiber.Router, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo)

	users := api.Group("/users", middleware.Auth(db))
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)
}
