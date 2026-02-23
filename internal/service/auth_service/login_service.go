package auth_service

import (
	"fmt"
	"golang-api/internal/config"
	"golang-api/internal/dto"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/internal/service"
	"golang-api/pkg/utils"
	"io"
	"net/http"
	"strings"
	"time"

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

	_, token, err := s.authenticate(data["email"].(string), data["password"].(string))
	if err != nil {
		if err.Error() == "invalid_credentials" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create token")
	}

	s.sendNotification(token)

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

func (s *loginService) sendNotification(token string) {
	mainURL := strings.TrimRight(config.MainUrl, "/")
	targetURL := fmt.Sprintf("%s/api/auth/login/notification", mainURL)

	req, err := http.NewRequest("POST", targetURL, nil)

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	io.ReadAll(resp.Body)
}
