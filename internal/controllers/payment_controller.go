/*
 * Project Name: controllers
 * File: payment_controller.go
 * Created Date: Saturday December 27th 2025
 * 
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 * 
 * Copyright (c) 2025-2026 Nova Ardiansyah, Org
 */

package controllers

import (
	"golang-api/internal/config"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"math"
	"path"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
)

const (
	PaymentTypeExpense    = 1
	PaymentTypeIncome     = 2
	PaymentTypeTransfer   = 3
	PaymentTypeWithdrawal = 4
)

type PaymentController struct {
	repo *repositories.PaymentRepository
}

type SummaryResponse struct {
	TotalBalance        int64           `json:"total_balance"`
	ScheduledExpense    int64           `json:"scheduled_expense"`
	TotalAfterScheduled int64           `json:"total_after_scheduled"`
	InitialBalance      int64           `json:"initial_balance"`
	Income              int64           `json:"income"`
	Expenses            int64           `json:"expenses"`
	Withdrawal          int64           `json:"withdrawal"`
	Transfer            int64           `json:"transfer"`
	Percents            SummaryPercents `json:"percents"`
	Period              SummaryPeriod   `json:"period"`
}

type SummaryPercents struct {
	Income     float64 `json:"income"`
	Expenses   float64 `json:"expenses"`
	Withdrawal float64 `json:"withdrawal"`
	Transfer   float64 `json:"transfer"`
}

type SummaryPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type AttachmentResponse struct {
	ID            int                 `json:"id"`
	URL           string              `json:"url"`
	Filepath      string              `json:"filepath"`
	Filename      string              `json:"filename"`
	Extension     string              `json:"extension"`
	FormattedSize string              `json:"formatted_size"`
	Original      *OriginalAttachment `json:"original"`
}

type OriginalAttachment struct {
	URL           string `json:"url"`
	FormattedSize string `json:"formatted_size"`
}

func NewPaymentController(repo *repositories.PaymentRepository) *PaymentController {
	return &PaymentController{repo: repo}
}

// Index godoc
// @Summary List payments
// @Description Get a paginated list of payments with optional filters
// @Tags payments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param type query int false "Type ID (1: Expense, 2: Income, 3: Transfer, 4: Withdrawal)"
// @Param account_id query int false "Account ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search query"
// @Success 200 {object} utils.PaginatedResponse{data=[]PaymentSwagger}
// @Failure 400 {object} utils.Response
// @Router /payments [get]
// @Security BearerAuth
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

// Show godoc
// @Summary Get payment details
// @Description Get detailed information about a specific payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} utils.Response{data=PaymentSwagger}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /payments/{id} [get]
// @Security BearerAuth
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

// Summary godoc
// @Summary Get payment summary
// @Description Get summary of payments within a date range
// @Tags payments
// @Accept json
// @Produce json
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} utils.Response{data=SummaryResponse}
// @Failure 422 {object} utils.ValidationErrorResponse
// @Router /payments/summary [get]
// @Security BearerAuth
func (ctrl *PaymentController) Summary(c *fiber.Ctx) error {
	data := make(map[string]interface{})

	rules := govalidator.MapData{
		"startDate": []string{"date"},
		"endDate":   []string{"date"},
	}

	errs := utils.ValidateJSON(c, &data, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	var startDate, endDate string

	if val, ok := data["startDate"].(string); ok && val != "" {
		startDate = val
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}

	if val, ok := data["endDate"].(string); ok && val != "" {
		endDate = val
	} else {
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}

	userID := c.Locals("user_id").(uint)
	db := config.GetDB()

	var totals struct {
		TotalIncome      int64
		TotalExpense     int64
		TotalWithdrawal  int64
		TotalTransfer    int64
		ScheduledExpense int64
	}

	db.Model(&models.Payment{}).
		Where("user_id = ?", userID).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Select(`
			SUM(CASE WHEN type_id = ? THEN amount ELSE 0 END) as total_income,
			SUM(CASE WHEN type_id = ? THEN amount ELSE 0 END) as total_expense,
			SUM(CASE WHEN type_id = ? THEN amount ELSE 0 END) as total_withdrawal,
			SUM(CASE WHEN type_id = ? THEN amount ELSE 0 END) as total_transfer,
			SUM(CASE WHEN type_id = ? AND is_scheduled = 1 THEN amount ELSE 0 END) as scheduled_expense
		`, PaymentTypeIncome, PaymentTypeExpense, PaymentTypeWithdrawal, PaymentTypeTransfer, PaymentTypeExpense).
		Scan(&totals)

	var totalBalance int64
	db.Model(&models.PaymentAccount{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(deposit), 0)").
		Scan(&totalBalance)

	initialBalance := totals.TotalIncome + totals.TotalExpense

	var percentIncome, percentExpense, percentWithdrawal, percentTransfer float64

	if initialBalance > 0 {
		percentIncome = math.Round((float64(totals.TotalIncome)/float64(initialBalance)*100)*100) / 100
		percentExpense = math.Round((float64(totals.TotalExpense)/float64(initialBalance)*100)*100) / 100
		percentWithdrawal = math.Round((float64(totals.TotalWithdrawal)/float64(initialBalance)*100)*100) / 100
		percentTransfer = math.Round((float64(totals.TotalTransfer)/float64(initialBalance)*100)*100) / 100
	}

	totalAfterScheduled := totalBalance - totals.ScheduledExpense

	response := SummaryResponse{
		TotalBalance:        totalBalance,
		ScheduledExpense:    totals.ScheduledExpense,
		TotalAfterScheduled: totalAfterScheduled,
		InitialBalance:      initialBalance,
		Income:              totals.TotalIncome,
		Expenses:            totals.TotalExpense,
		Withdrawal:          totals.TotalWithdrawal,
		Transfer:            totals.TotalTransfer,
		Percents: SummaryPercents{
			Income:     percentIncome,
			Expenses:   percentExpense,
			Withdrawal: percentWithdrawal,
			Transfer:   percentTransfer,
		},
		Period: SummaryPeriod{
			StartDate: startDate,
			EndDate:   endDate,
		},
	}

	return utils.SuccessResponse(c, "Summary retrieved successfully", response)
}

// GetAttachments godoc
// @Summary Get payment attachments
// @Description Get a list of attachments for a specific payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} utils.Response{data=[]AttachmentResponse}
// @Failure 400 {object} utils.Response
// @Router /payments/{id}/attachments [get]
// @Security BearerAuth
func (ctrl *PaymentController) GetAttachments(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid payment ID")
	}

	payment, err := ctrl.repo.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Payment not found")
	}

	var attachments []string
	if len(payment.Attachments) > 0 {
		if err := json.Unmarshal(payment.Attachments, &attachments); err != nil {
			attachments = []string{}
		}
	}

	var responseData []AttachmentResponse
	for i, attachment := range attachments {
		filename := path.Base(attachment)
		ext := strings.TrimPrefix(path.Ext(filename), ".")
		nameOnly := strings.TrimSuffix(filename, path.Ext(filename))

		mediumName := "medium-" + nameOnly + "." + ext
		mediumPath := "images/payment/" + mediumName
		originalPath := "images/payment/" + filename

		mediumURL := config.CdnUrl + "/" + mediumPath
		originalURL := config.CdnUrl + "/" + originalPath

		responseData = append(responseData, AttachmentResponse{
			ID:            i + 1,
			URL:           mediumURL,
			Filepath:      mediumPath,
			Filename:      mediumName,
			Extension:     ext,
			FormattedSize: "0 KB",
			Original: &OriginalAttachment{
				URL:           originalURL,
				FormattedSize: "0 KB",
			},
		})
	}

	message := "No attachments found"
	if len(responseData) > 0 {
		message = "Attachments found"
	}

	return utils.SuccessResponse(c, message, responseData)
}
