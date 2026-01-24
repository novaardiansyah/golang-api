/*
 * Project Name: routes
 * File: routes.go
 * Created Date: Saturday December 27th 2025
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2025-2026 Nova Ardiansyah, Org
 */

package routes

import (
	"golang-api/internal/config"
	"golang-api/internal/controllers"
	"golang-api/internal/middleware"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	app.Use(middleware.GlobalLimiter())
	app.Static("/", "./public")

	userRepo := repositories.NewUserRepository(db)

	userController := controllers.NewUserController(userRepo)
	authController := controllers.NewAuthController(db)

	api := app.Group("/api")

	api.Get("/documentation/*", swagger.HandlerDefault)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is running",
		})
	})

	api.Get("/test-email", middleware.Auth(db), func(c *fiber.Ctx) error {
		to := c.Query("to", "admin@novadev.my.id")

		err := utils.SendEmail(to, "Testing Golang Email", map[string]any{
			"Name":    "User Testing",
			"Message": "This is a dynamic message from the route!",
		}, "resources/views/emails/main.html", "resources/views/emails/test_content.html")

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to send email: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Email has been sent to " + to,
		})
	})

	auth := api.Group("/auth")
	auth.Use(middleware.AuthLimiter())

	auth.Post("/login", authController.Login)
	auth.Get("/validate-token", middleware.Auth(db), authController.ValidateToken)
	auth.Post("/logout", middleware.Auth(db), authController.Logout)
	auth.Post("/change-password", middleware.Auth(db), authController.ChangePassword)

	users := api.Group("/users", middleware.Auth(db))
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)

	paymentRepo := repositories.NewPaymentRepository(db)
	paymentController := controllers.NewPaymentController(paymentRepo)

	payments := api.Group("/payments", middleware.Auth(db))
	payments.Get("/", paymentController.Index)
	payments.Get("/summary", paymentController.Summary)
	payments.Get("/:id", paymentController.Show)

	notificationController := controllers.NewNotificationController()

	notifications := api.Group("/notifications", middleware.Auth(db))
	notifications.Put("/settings", notificationController.UpdateSettings)

	fileDownloadRepo := repositories.NewFileDownloadRepository(db)
	fileDownloadController := controllers.NewFileDownloadController(fileDownloadRepo)

	files := api.Group("/files/d", middleware.Auth(db))
	files.Get("/:uid", fileDownloadController.GetFiles)

	paymentGoalRepo := repositories.NewPaymentGoalRepository(db)
	paymentGoalController := controllers.NewPaymentGoalController(paymentGoalRepo)

	paymentGoals := api.Group("/payment-goals", middleware.Auth(db))
	paymentGoals.Get("/", paymentGoalController.Index)
	paymentGoals.Get("/overview", paymentGoalController.Overview)
	paymentGoals.Get("/:id", paymentGoalController.Show)

	paymentAccountRepo := repositories.NewPaymentAccountRepository(db)
	paymentAccountController := controllers.NewPaymentAccountController(paymentAccountRepo)

	paymentAccounts := api.Group("/payment-accounts", middleware.Auth(db))
	paymentAccounts.Get("/", paymentAccountController.Index)
}
