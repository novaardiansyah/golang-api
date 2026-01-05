package controllers

import (
	"golang-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type EmailController struct {
	repo *repositories.EmailRepository
}

func NewEmailController(repo *repositories.EmailRepository) *EmailController {
	return &EmailController{repo: repo}
}

func (ctrl *EmailController) Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}
