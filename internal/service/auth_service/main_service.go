package auth_service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MainService interface {
	Login(c *fiber.Ctx) error
}

type mainService struct {
	loginService LoginService
}

func NewMainService(db *gorm.DB) MainService {
	return &mainService{
		loginService: NewLoginService(db),
	}
}

func (s *mainService) Login(c *fiber.Ctx) error {
	return s.loginService.Login(c)
}
