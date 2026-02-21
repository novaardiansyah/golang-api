/*
 * Project Name: controllers
 * File: payment_type_controller.go
 * Created Date: Saturday February 21st 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2025-2026 Nova Ardiansyah, Org
 */

package controllers

import (
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentTypeController struct {
	repo *repositories.PaymentTypeRepository
}

func NewPaymentTypeController(repo *repositories.PaymentTypeRepository) *PaymentTypeController {
	return &PaymentTypeController{repo: repo}
}

// Index godoc
// @Summary List payment types
// @Description Get a paginated list of payment types
// @Tags payment_types
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]PaymentTypeSwagger}
// @Failure 401 {object} utils.UnauthorizedResponse
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 500 {object} utils.SimpleErrorResponse
// @Router /payment-types [get]
// @Security BearerAuth
func (ctrl *PaymentTypeController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	} else if perPage > 100 {
		perPage = 100
	}

	total, err := ctrl.repo.Count()

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count payment types")
	}

	paymentTypes, err := ctrl.repo.FindAllPaginated(page, perPage)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve payment types")
	}

	return utils.PaginatedSuccessResponse(c, "Payment types retrieved successfully", paymentTypes, page, perPage, total, len(paymentTypes))
}
