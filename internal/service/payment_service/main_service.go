package payment_service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MainService interface {
	Store(c *fiber.Ctx) error
}

type mainService struct {
	storeService StoreService
}

func NewMainService(db *gorm.DB) MainService {
	return &mainService{
		storeService: NewStoreService(db),
	}
}

func (s *mainService) Store(c *fiber.Ctx) error {
	return s.storeService.Store(c)
}
