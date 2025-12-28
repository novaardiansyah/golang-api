package controllers

import (
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	repo *repositories.PaymentRepository
}

func NewPaymentController(repo *repositories.PaymentRepository) *PaymentController {
	return &PaymentController{repo: repo}
}

func (ctrl *PaymentController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	typeID, _ := strconv.Atoi(c.Query("type", "0"))
	accountID, _ := strconv.Atoi(c.Query("account_id", "0"))

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}

	filter := repositories.PaymentFilter{
		DateFrom:  c.Query("date_from"),
		DateTo:    c.Query("date_to"),
		Type:      typeID,
		AccountID: accountID,
		Search:    c.Query("search"),
	}

	total, err := ctrl.repo.Count(filter)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count payments")
	}

	payments, err := ctrl.repo.FindAllPaginated(page, perPage, filter)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve payments")
	}

	return utils.PaginatedSuccessResponse(c, "Payments retrieved successfully", payments, page, perPage, total, len(payments))
}

func (ctrl *PaymentController) Show(c *fiber.Ctx) error {
	paymentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid payment ID")
	}

	payment, err := ctrl.repo.FindByID(paymentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Payment not found")
	}

	return utils.SuccessResponse(c, "Payment retrieved successfully", payment)
}
