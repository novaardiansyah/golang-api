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

type UpdateNotificationSettingsRequest struct {
	HasAllowNotification int    `json:"has_allow_notification" validate:"numeric_between:0,1"`
	NotificationToken    string `json:"notification_token" validate:"max:255"`
}

func NewNotificationController() *NotificationController {
	return &NotificationController{}
}

// UpdateSettings godoc
// @Summary Update notification settings
// @Description Update user notification allowance and Expo push token
// @Tags notifications
// @Accept json
// @Produce json
// @Param settings body UpdateNotificationSettingsRequest true "Notification settings"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.ValidationErrorResponse
// @Failure 500 {object} utils.Response
// @Router /notifications/settings [put]
// @Security BearerAuth
func (ctrl *NotificationController) UpdateSettings(c *fiber.Ctx) error {
	data := make(map[string]interface{})

	rules := govalidator.MapData{
		"has_allow_notification": []string{"numeric_between:0,1"},
		"notification_token":     []string{"max:255"},
	}

	errs := utils.ValidateJSON(c, &data, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	notificationToken := ""
	if val, ok := data["notification_token"].(string); ok {
		notificationToken = val
	}

	if notificationToken != "" {
		if !validateExpoToken(notificationToken) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Expo push token format")
		}
	}

	userID := c.Locals("user_id").(uint)
	db := config.GetDB()

	updates := make(map[string]interface{})

	if val, ok := data["has_allow_notification"]; ok {
		updates["has_allow_notification"] = val
	}

	if notificationToken != "" {
		updates["notification_token"] = notificationToken
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
