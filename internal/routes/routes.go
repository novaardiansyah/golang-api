package routes

import (
	"golang-api/internal/config"
	"golang-api/internal/middleware"
	"golang-api/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	app.Use(middleware.GlobalLimiter())
	app.Static("/", "./public")

	api := app.Group("/api")

	api.Get("/documentation/*", swagger.HandlerDefault)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is running",
			"data": map[string]any{
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
				"timezone":  "Asia/Jakarta",
			},
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

	AuthRoutes(api, db)
	UserRoutes(api, db)
	PaymentRoutes(api, db)
	NotificationRoutes(api, db)
	FileRoutes(api, db)
	PaymentGoalRoutes(api, db)
	PaymentAccountRoutes(api, db)
	ActivityLogRoutes(api, db)
	UptimeMonitorRoutes(api, db)
	UptimeMonitorLogRoutes(api, db)
}
