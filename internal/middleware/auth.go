package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"golang-api/internal/config"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(db *gorm.DB) fiber.Handler {
  userRepo := repositories.NewUserRepository(db)

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: No token provided",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token format",
			})
		}

		parts := strings.SplitN(tokenString, "|", 2)
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token format",
			})
		}

		tokenID := parts[0]
		plainTextToken := parts[1]

		hash := sha256.Sum256([]byte(plainTextToken))
		hashedToken := hex.EncodeToString(hash[:])

		db := config.GetDB()
		var token models.PersonalAccessToken

		result := db.Where("id = ? AND token = ?", tokenID, hashedToken).First(&token)
		if result.Error != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token",
			})
		}

		if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Token expired",
			})
		}

		db.Model(&token).Update("last_used_at", time.Now())

    UserId := token.TokenableID

    user, err := userRepo.FindByID(UserId)

    if err != nil {
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "success": false,
        "message": "Unauthorized: User not found",
      })
    }

		c.Locals("token", token)
		c.Locals("user_id", UserId)
    c.Locals("user", *user)

		return c.Next()
	}
}
