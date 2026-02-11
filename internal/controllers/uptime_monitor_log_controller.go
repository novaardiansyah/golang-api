/*
 * Project Name: controllers
 * File: uptime_monitor_log_controller.go
 * Created Date: Wednesday February 11th 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package controllers

import (
	"golang-api/internal/dto"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
	"gorm.io/gorm"
)

type UptimeMonitorLogController struct {
	repo *repositories.UptimeMonitorLogRepository
}

func NewUptimeMonitorLogController(db *gorm.DB) *UptimeMonitorLogController {
	return &UptimeMonitorLogController{
		repo: repositories.NewUptimeMonitorLogRepository(db),
	}
}

// Index godoc
// @Summary List uptime monitor logs
// @Description Get a paginated list of uptime monitor logs
// @Tags uptime_monitor_logs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param monitor_id query int false "Filter by monitor ID"
// @Success 200 {object} utils.PaginatedResponse{data=[]UptimeMonitorLogSwagger}
// @Router /uptime-monitor-logs [get]
// @Security BearerAuth
func (ctrl *UptimeMonitorLogController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	monitorID, _ := strconv.Atoi(c.Query("monitor_id", "0"))

	if page < 1 {
		page = 1
	}

	total, err := ctrl.repo.Count(uint(monitorID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count monitor logs")
	}

	logs, err := ctrl.repo.FindAllPaginated(page, perPage, uint(monitorID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve monitor logs")
	}

	return utils.PaginatedSuccessResponse(c, "Monitor logs retrieved successfully", logs, page, perPage, total, len(logs))
}

// Show godoc
// @Summary Get uptime monitor log details
// @Description Get a single uptime monitor log by ID
// @Tags uptime_monitor_logs
// @Accept json
// @Produce json
// @Param id path int true "Log ID"
// @Success 200 {object} utils.Response{data=UptimeMonitorLogSwagger}
// @Router /uptime-monitor-logs/{id} [get]
// @Security BearerAuth
func (ctrl *UptimeMonitorLogController) Show(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	log, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Monitor log not found")
	}

	return utils.SuccessResponse(c, "Monitor log retrieved successfully", log)
}

// Store godoc
// @Summary Store a new uptime monitor log
// @Description Store a new uptime monitor log
// @Tags uptime_monitor_logs
// @Accept json
// @Produce json
// @Param log body dto.StoreUptimeMonitorLogRequest true "Log data"
// @Success 201 {object} utils.Response{data=UptimeMonitorLogSwagger}
// @Router /uptime-monitor-logs [post]
// @Security BearerAuth
func (ctrl *UptimeMonitorLogController) Store(c *fiber.Ctx) error {
	var request dto.StoreUptimeMonitorLogRequest

	rules := govalidator.MapData{
		"uptime_monitor_id": []string{"required", "numeric"},
		"status_code":       []string{"required", "numeric"},
		"response_time_ms":  []string{"required", "numeric"},
		"checked_at":        []string{"required"},
	}

	errs := utils.ValidateJSON(c, &request, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	log := models.UptimeMonitorLog{
		UptimeMonitorID: request.UptimeMonitorID,
		StatusCode:      request.StatusCode,
		ResponseTimeMs:  request.ResponseTimeMs,
		IsHealthy:       request.IsHealthy,
		ErrorMessage:    request.ErrorMessage,
		CheckedAt:       request.CheckedAt,
	}

	if err := ctrl.repo.Store(&log); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to store monitor log")
	}

	return utils.SuccessResponse(c, "Monitor log stored successfully", log)
}

// Destroy godoc
// @Summary Delete an uptime monitor log
// @Description Delete an existing uptime monitor log
// @Tags uptime_monitor_logs
// @Accept json
// @Produce json
// @Param id path int true "Log ID"
// @Success 200 {object} utils.Response
// @Router /uptime-monitor-logs/{id} [delete]
// @Security BearerAuth
func (ctrl *UptimeMonitorLogController) Destroy(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := ctrl.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete monitor log")
	}

	return utils.SuccessResponse(c, "Monitor log deleted successfully", nil)
}
