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

  if page < 1 {
    page = 1
  }

  if perPage < 1 {
    perPage = 10
  }

  total, err := ctrl.repo.Count()

  if err != nil {
    return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count payments")
  }

  payments, err := ctrl.repo.FindAllPaginated(page, perPage)

  if err != nil {
    return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve payments")
  }

  return utils.PaginatedSuccessResponse(c, "Payments retrieved successfully", payments, page, perPage, total, len(payments))
}