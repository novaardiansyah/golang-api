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

// Index godoc
// @Summary List payment goals
// @Description Get a list of payment goals for the authenticated user
// @Tags payment-goals
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]PaymentGoalSwagger}
// @Failure 400 {object} utils.Response
// @Router /payment-goals [get]
// @Security BearerAuth
func (ctrl *PaymentGoalController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	goals, err := ctrl.repo.FindAll(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve payment goals")
	}

	return utils.SuccessResponse(c, "Payment goals retrieved successfully", goals)
}

// Show godoc
// @Summary Get payment goal details
// @Description Get detailed information about a specific payment goal
// @Tags payment-goals
// @Accept json
// @Produce json
// @Param id path int true "Payment Goal ID"
// @Success 200 {object} utils.Response{data=PaymentGoalSwagger}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /payment-goals/{id} [get]
// @Security BearerAuth
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

type OverviewResponse struct {
	TotalGoals  int64  `json:"total_goals"`
	Completed   int64  `json:"completed"`
	SuccessRate string `json:"success_rate"`
}

// Overview godoc
// @Summary Get payment goals overview
// @Description Get overview statistics of payment goals for the authenticated user
// @Tags payment-goals
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=OverviewResponse}
// @Failure 400 {object} utils.Response
// @Router /payment-goals/overview [get]
// @Security BearerAuth
func (ctrl *PaymentGoalController) Overview(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	totalGoals, completedGoals, err := ctrl.repo.GetOverview(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve overview")
	}

	successRate := "0%"
	if totalGoals > 0 {
		rate := float64(completedGoals) / float64(totalGoals) * 100
		successRate = utils.FormatPercent(int(rate))
	}

	response := OverviewResponse{
		TotalGoals:  totalGoals,
		Completed:   completedGoals,
		SuccessRate: successRate,
	}

	return utils.SuccessResponse(c, "Overview retrieved successfully", response)
}
