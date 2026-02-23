package auth_service

import (
	"encoding/json"
	"fmt"
	"golang-api/internal/dto"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/internal/service"
	"golang-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
	"gorm.io/gorm"
)

type LoginService interface {
	Login(c *fiber.Ctx) error
}

type loginService struct {
	userRepo        *repositories.UserRepository
	activityLogRepo *repositories.ActivityLogRepository
	authService     service.AuthService
}

func NewLoginService(db *gorm.DB) LoginService {
	return &loginService{
		userRepo:        repositories.NewUserRepository(db),
		activityLogRepo: repositories.NewActivityLogRepository(db),
		authService:     service.NewAuthService(db),
	}
}

func (s *loginService) Login(c *fiber.Ctx) error {
	data, errs := s.validate(c)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	user, token, err := s.authenticate(data["email"].(string), data["password"].(string))
	if err != nil {
		if err.Error() == "invalid_credentials" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create token")
	}

	s.createActivityLog(c, user)

	return utils.SuccessResponse(c, "Login successful", dto.LoginResponse{
		Token: token,
	})
}

func (s *loginService) validate(c *fiber.Ctx) (map[string]interface{}, map[string][]string) {
	data := make(map[string]interface{})

	rules := govalidator.MapData{
		"email":    []string{"required", "email"},
		"password": []string{"required", "min:6"},
	}

	errs := utils.ValidateJSON(c, &data, rules)
	if errs != nil {
		return nil, errs
	}

	return data, nil
}

func (s *loginService) authenticate(email, password string) (*models.User, string, error) {
	return s.authService.Login(email, password)
}

func (s *loginService) createActivityLog(c *fiber.Ctx, user *models.User) {
	properties, _ := json.Marshal(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})

	activityLog := models.ActivityLog{
		LogName:        "Resource",
		Description:    fmt.Sprintf("User %s has successfully authenticated via the API service", user.Name),
		SubjectID:      &user.ID,
		SubjectType:    utils.String("App\\Models\\User"),
		Event:          "Login",
		CauserID:       user.ID,
		CauserType:     "App\\Models\\User",
		PrevProperties: utils.RawMessage(json.RawMessage("[]")),
		Properties:     properties,
		IPAddress:      utils.String(c.IP()),
		UserAgent:      utils.String(string(c.Request().Header.UserAgent())),
		Referer:        utils.String(c.Get("Referer")),
	}

	s.activityLogRepo.Store(&activityLog)
}
