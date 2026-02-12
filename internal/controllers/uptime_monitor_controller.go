/*
 * Project Name: controllers
 * File: uptime_monitor_controller.go
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
	"golang-api/internal/service"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
	"gorm.io/gorm"
)

type UptimeMonitorController struct {
	repo    *repositories.UptimeMonitorRepository
	service *service.UptimeMonitorService
}

func NewUptimeMonitorController(db *gorm.DB) *UptimeMonitorController {
	return &UptimeMonitorController{
		repo:    repositories.NewUptimeMonitorRepository(db),
		service: service.NewUptimeMonitorService(db),
	}
}

// RunChecks godoc
// @Summary Run uptime monitor checks
// @Description Run uptime monitor checks for all active monitors that are due
// @Tags uptime_monitors
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /uptime-monitors/run-checks [post]
// @Security BearerAuth
func (ctrl *UptimeMonitorController) RunChecks(c *fiber.Ctx) error {
	results := ctrl.service.RunScheduledChecks()
	return utils.SuccessResponse(c, "Uptime monitor checks completed", results)
}

// Index godoc
// @Summary List uptime monitors
// @Description Get a paginated list of uptime monitors
// @Tags uptime_monitors
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by name or code"
// @Success 200 {object} utils.PaginatedResponse{data=[]UptimeMonitorSwagger}
// @Router /uptime-monitors [get]
// @Security BearerAuth
func (ctrl *UptimeMonitorController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}

	total, err := ctrl.repo.Count(search)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count uptime monitors")
	}

	monitors, err := ctrl.repo.FindAllPaginated(page, perPage, search)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve uptime monitors")
	}

	return utils.PaginatedSuccessResponse(c, "Uptime monitors retrieved successfully", monitors, page, perPage, total, len(monitors))
}

// Show godoc
// @Summary Get uptime monitor details
// @Description Get a single uptime monitor by ID
// @Tags uptime_monitors
// @Accept json
// @Produce json
// @Param id path int true "Monitor ID"
// @Success 200 {object} utils.Response{data=UptimeMonitorSwagger}
// @Router /uptime-monitors/{id} [get]
// @Security BearerAuth
func (ctrl *UptimeMonitorController) Show(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	monitor, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Uptime monitor not found")
	}

	return utils.SuccessResponse(c, "Uptime monitor retrieved successfully", monitor)
}

// Store godoc
// @Summary Store a new uptime monitor
// @Description Store a new uptime monitor
// @Tags uptime_monitors
// @Accept json
// @Produce json
// @Param monitor body dto.StoreUptimeMonitorRequest true "Monitor data"
// @Success 201 {object} utils.Response{data=UptimeMonitorSwagger}
// @Router /uptime-monitors [post]
// @Security BearerAuth
func (ctrl *UptimeMonitorController) Store(c *fiber.Ctx) error {
	var request dto.StoreUptimeMonitorRequest

	rules := govalidator.MapData{
		"code":     []string{"required"},
		"url":      []string{"required", "url"},
		"name":     []string{"required"},
		"interval": []string{"required", "numeric"},
	}

	errs := utils.ValidateJSON(c, &request, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	monitor := models.UptimeMonitor{
		Code:     request.Code,
		URL:      request.URL,
		Name:     request.Name,
		Interval: request.Interval,
		IsActive: request.IsActive,
	}

	if err := ctrl.repo.Store(&monitor); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to store uptime monitor")
	}

	return utils.SuccessResponse(c, "Uptime monitor stored successfully", monitor)
}

// Update godoc
// @Summary Update an uptime monitor
// @Description Update an existing uptime monitor
// @Tags uptime_monitors
// @Accept json
// @Produce json
// @Param id path int true "Monitor ID"
// @Param monitor body dto.UpdateUptimeMonitorRequest true "Monitor data"
// @Success 200 {object} utils.Response{data=UptimeMonitorSwagger}
// @Router /uptime-monitors/{id} [put]
// @Security BearerAuth
func (ctrl *UptimeMonitorController) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	monitor, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Uptime monitor not found")
	}

	var request dto.UpdateUptimeMonitorRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if request.Code != "" {
		monitor.Code = request.Code
	}
	if request.URL != "" {
		monitor.URL = request.URL
	}
	if request.Name != "" {
		monitor.Name = request.Name
	}
	if request.Interval != 0 {
		monitor.Interval = request.Interval
	}
	monitor.IsActive = request.IsActive

	if err := ctrl.repo.Update(monitor); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update uptime monitor")
	}

	return utils.SuccessResponse(c, "Uptime monitor updated successfully", monitor)
}

// Destroy godoc
// @Summary Delete an uptime monitor
// @Description Delete an existing uptime monitor
// @Tags uptime_monitors
// @Accept json
// @Produce json
// @Param id path int true "Monitor ID"
// @Success 200 {object} utils.Response
// @Router /uptime-monitors/{id} [delete]
// @Security BearerAuth
func (ctrl *UptimeMonitorController) Destroy(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := ctrl.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete uptime monitor")
	}

	return utils.SuccessResponse(c, "Uptime monitor deleted successfully", nil)
}
