package controllers

import (
	"golang-api/internal/config"
	"golang-api/internal/models"
	"golang-api/pkg/utils"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
)

type NotificationController struct{}

func NewNotificationController() *NotificationController {
	return &NotificationController{}
}

type UpdateNotificationSettingsRequest struct {
	HasAllowNotification *bool  `json:"has_allow_notification"`
	NotificationToken    string `json:"notification_token"`
}

func (ctrl *NotificationController) UpdateSettings(c *fiber.Ctx) error {
	var req UpdateNotificationSettingsRequest

	rules := govalidator.MapData{
		"has_allow_notification": []string{"bool"},
		"notification_token":     []string{"max:255"},
	}

	errs := utils.ValidateJSONStruct(c, &req, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	if req.NotificationToken != "" {
		if !validateExpoToken(req.NotificationToken) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Expo push token format")
		}
	}

	userID := c.Locals("user_id").(uint)
	db := config.GetDB()

	updates := make(map[string]interface{})
	if req.HasAllowNotification != nil {
		updates["has_allow_notification"] = *req.HasAllowNotification
	}
	updates["notification_token"] = req.NotificationToken

	if err := db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update notification settings")
	}

	return utils.SuccessResponse(c, "Notification settings updated successfully", nil)
}

func validateExpoToken(token string) bool {
	expoPattern := regexp.MustCompile(`^ExponentPushToken\[([a-zA-Z0-9\-_]+)\]$`)
	hexPattern := regexp.MustCompile(`^[a-f0-9]{32}$`)
	return expoPattern.MatchString(token) || hexPattern.MatchString(token)
}
