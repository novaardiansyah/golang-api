package controllers

import (
	"golang-api/internal/config"
	"golang-api/internal/models"
	"golang-api/pkg/utils"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

type NotificationController struct{}

func NewNotificationController() *NotificationController {
	return &NotificationController{}
}

func (ctrl *NotificationController) UpdateSettings(c *fiber.Ctx) error {
	rules := map[string]utils.FieldRule{
		"has_allow_notification": {Type: "bool"},
		"notification_token":     {Type: "string", Max: 255},
	}

	data, errs := utils.ValidateJSONMap(c, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	hasAllowNotification := utils.GetBool(data, "has_allow_notification")
	notificationToken := utils.GetString(data, "notification_token")

	if notificationToken != nil && *notificationToken != "" {
		if !validateExpoToken(*notificationToken) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Expo push token format")
		}
	}

	userID := c.Locals("user_id").(uint)
	db := config.GetDB()

	updates := make(map[string]interface{})
	if hasAllowNotification != nil {
		updates["has_allow_notification"] = *hasAllowNotification
	}
	if notificationToken != nil {
		updates["notification_token"] = *notificationToken
	}

	if len(updates) > 0 {
		if err := db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update notification settings")
		}
	}

	return utils.SuccessResponse(c, "Notification settings updated successfully", nil)
}

func validateExpoToken(token string) bool {
	expoPattern := regexp.MustCompile(`^ExponentPushToken\[([a-zA-Z0-9\-_]+)\]$`)
	hexPattern := regexp.MustCompile(`^[a-f0-9]{32}$`)
	return expoPattern.MatchString(token) || hexPattern.MatchString(token)
}
