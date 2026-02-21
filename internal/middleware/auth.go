package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(db *gorm.DB) fiber.Handler {
	PersonalAccessTokenRepo := repositories.NewPersonalAccessTokenRepository(db)

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

		tokenIDStr := parts[0]
		plainTextToken := parts[1]

		tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token format",
			})
		}

		hash := sha256.Sum256([]byte(plainTextToken))
		hashedToken := hex.EncodeToString(hash[:])

		result, err := PersonalAccessTokenRepo.FindByIDAndHashedTokenWithUser(tokenID, hashedToken)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Invalid token",
			})
		}

		if result.ExpiresAt != nil && result.ExpiresAt.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized: Token expired",
			})
		}

		token := &models.PersonalAccessToken{
			ID:            result.ID,
			TokenableType: result.TokenableType,
			TokenableID:   result.TokenableID,
			Name:          result.Name,
			Token:         result.Token,
			Abilities:     result.Abilities,
			LastUsedAt:    result.LastUsedAt,
			ExpiresAt:     result.ExpiresAt,
			CreatedAt:     result.CreatedAt,
			UpdatedAt:     result.UpdatedAt,
		}

		fields := map[string]interface{}{"last_used_at": time.Now()}
		PersonalAccessTokenRepo.UpdateFields(token, fields)

		UserId := token.TokenableID

		c.Locals("token", *token)
		c.Locals("user_id", UserId)
		c.Locals("user_name", result.UserName)

		return c.Next()
	}
}
