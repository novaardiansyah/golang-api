/*
 * Project Name: controllers
 * File: payment_account.go
 * Created Date: Thursday January 22nd 2026
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

type PaymentAccountController struct {
	repo *repositories.PaymentAccountRepository
}

func NewPaymentAccountController(repo *repositories.PaymentAccountRepository) *PaymentAccountController {
	return &PaymentAccountController{repo: repo}
}

// Index godoc
// @Summary List payment accounts
// @Description Get a paginated list of payment accounts
// @Tags payment_accounts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]PaymentAccountSwagger}
// @Failure 400 {object} utils.Response
// @Router /payment-accounts [get]
// @Security BearerAuth
func (ctrl *PaymentAccountController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
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

	total, err := ctrl.repo.Count(userID)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count payment accounts")
	}

	paymentAccounts, err := ctrl.repo.FindAllPaginated(userID, page, perPage)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve payment accounts")
	}

	return utils.PaginatedSuccessResponse(c, "Payment accounts retrieved successfully", paymentAccounts, page, perPage, total, len(paymentAccounts))
}
