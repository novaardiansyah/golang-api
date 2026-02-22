package payment_service

import (
	"encoding/json"
	"fmt"
	"golang-api/internal/config"
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

type generateReportRequest struct {
	ReportType string `json:"report_type"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Periode    string `json:"periode"`
}

func (s *generateReportService) GenerateReport(c *fiber.Ctx) error {
	var payload generateReportRequest

	rules := govalidator.MapData{
		"report_type": []string{"required", "in:daily,monthly,date_range"},
	}

	errs := utils.ValidateJSON(c, &payload, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	validationErrs := make(map[string][]string)

	if payload.ReportType == "date_range" {
		dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

		if payload.StartDate == "" {
			validationErrs["start_date"] = []string{"The start date is required for custom date range report."}
		} else if !dateRegex.MatchString(payload.StartDate) {
			validationErrs["start_date"] = []string{"The start date must be in format Y-m-d."}
		}

		if payload.EndDate == "" {
			validationErrs["end_date"] = []string{"The end date is required for custom date range report."}
		} else if !dateRegex.MatchString(payload.EndDate) {
			validationErrs["end_date"] = []string{"The end date must be in format Y-m-d."}
		}

		if len(validationErrs) == 0 {
			start, _ := time.Parse("2006-01-02", payload.StartDate)
			end, _ := time.Parse("2006-01-02", payload.EndDate)
			if end.Before(start) {
				validationErrs["end_date"] = []string{"The end date must be after or equal to start date."}
			}
		}
	}

	if payload.ReportType == "monthly" {
		monthRegex := regexp.MustCompile(`^\d{4}-\d{2}$`)

		if payload.Periode == "" {
			validationErrs["periode"] = []string{"The periode (month) is required for monthly report."}
		} else if !monthRegex.MatchString(payload.Periode) {
			validationErrs["periode"] = []string{"The periode must be in format Y-m."}
		}
	}

	if len(validationErrs) > 0 {
		return utils.ValidationError(c, validationErrs)
	}

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
