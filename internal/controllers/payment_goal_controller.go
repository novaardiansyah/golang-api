package controllers

import (
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentGoalController struct {
	repo *repositories.PaymentGoalRepository
}

func NewPaymentGoalController(repo *repositories.PaymentGoalRepository) *PaymentGoalController {
	return &PaymentGoalController{repo: repo}
}

func (ctrl *PaymentGoalController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	goals, err := ctrl.repo.FindAll(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve payment goals")
	}

	return utils.SuccessResponse(c, "Payment goals retrieved successfully", goals)
}

func (ctrl *PaymentGoalController) Show(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	goalID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid payment goal ID")
	}

	goal, err := ctrl.repo.FindByID(goalID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Payment goal not found")
	}

	return utils.SuccessResponse(c, "Payment goal retrieved successfully", goal)
}
