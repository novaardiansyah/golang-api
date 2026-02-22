package payment_service

import (
	"encoding/json"
	"fmt"
	"golang-api/internal/config"
	"golang-api/internal/dto"
	"golang-api/pkg/utils"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
)

type GenerateReportService interface {
	GenerateReport(c *fiber.Ctx) error
}

type generateReportService struct{}

func NewGenerateReportService() GenerateReportService {
	return &generateReportService{}
}

func (s *generateReportService) GenerateReport(c *fiber.Ctx) error {
	var payload dto.GenerateReportRequest

	validateErrors := s.validate(c, &payload)
	if validateErrors != nil {
		return utils.ValidationError(c, validateErrors)
	}

	return s.forwardRequest(c, &payload)
}

func (s *generateReportService) validate(c *fiber.Ctx, payload *dto.GenerateReportRequest) map[string][]string {
	rules := govalidator.MapData{
		"report_type": []string{"required", "in:daily,monthly,date_range"},
	}

	errs := utils.ValidateJSON(c, payload, rules)
	if errs != nil {
		return errs
	}

	validationErrs := make(map[string][]string)

	if payload.ReportType == "date_range" {
		s.validateDateRange(payload, validationErrs)
	}

	if payload.ReportType == "monthly" {
		s.validateMonthly(payload, validationErrs)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (s *generateReportService) validateDateRange(payload *dto.GenerateReportRequest, errs map[string][]string) {
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	if payload.StartDate == "" {
		errs["start_date"] = []string{"The start date is required for custom date range report."}
	} else if !dateRegex.MatchString(payload.StartDate) {
		errs["start_date"] = []string{"The start date must be in format Y-m-d."}
	}

	if payload.EndDate == "" {
		errs["end_date"] = []string{"The end date is required for custom date range report."}
	} else if !dateRegex.MatchString(payload.EndDate) {
		errs["end_date"] = []string{"The end date must be in format Y-m-d."}
	}

	if len(errs) == 0 {
		start, _ := time.Parse("2006-01-02", payload.StartDate)
		end, _ := time.Parse("2006-01-02", payload.EndDate)
		if end.Before(start) {
			errs["end_date"] = []string{"The end date must be after or equal to start date."}
		}
	}
}

func (s *generateReportService) validateMonthly(payload *dto.GenerateReportRequest, errs map[string][]string) {
	monthRegex := regexp.MustCompile(`^\d{4}-\d{2}$`)

	if payload.Periode == "" {
		errs["periode"] = []string{"The periode (month) is required for monthly report."}
	} else if !monthRegex.MatchString(payload.Periode) {
		errs["periode"] = []string{"The periode must be in format Y-m."}
	}
}

func (s *generateReportService) forwardRequest(c *fiber.Ctx, payload *dto.GenerateReportRequest) error {
	mainURL := strings.TrimRight(config.MainUrl, "/")
	targetURL := fmt.Sprintf("%s/api/payments/generate-report", mainURL)

	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	authHeader := c.Get("Authorization")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadGateway, "Failed to connect to main server")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read response from main server")
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	return c.Status(resp.StatusCode).Send(respBody)
}
