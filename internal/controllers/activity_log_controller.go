/*
 * Project Name: controllers
 * File: activity_log_controller.go
 * Created Date: Tuesday February 10th 2026
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

type ActivityLogController struct {
	repo *repositories.ActivityLogRepository
}

func NewActivityLogController(db *gorm.DB) *ActivityLogController {
	return &ActivityLogController{
		repo: repositories.NewActivityLogRepository(db),
	}
}

// Index godoc
// @Summary List activity logs
// @Description Get a paginated list of activity logs
// @Tags activity_logs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]ActivityLogSwagger}
// @Failure 400 {object} utils.Response
// @Router /activity-logs [get]
// @Security BearerAuth
func (ctrl *ActivityLogController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	} else if perPage > 100 {
		perPage = 100
	}

	total, err := ctrl.repo.Count(search)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count activity logs")
	}

	activityLogs, err := ctrl.repo.FindAllPaginated(page, perPage, search)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve activity logs")
	}

	return utils.PaginatedSuccessResponse(c, "Activity logs retrieved successfully", activityLogs, page, perPage, total, len(activityLogs))
}

// Store godoc
// @Summary Store a new activity log
// @Description Store a new activity log
// @Tags activity_logs
// @Accept json
// @Produce json
// @Param activity_log body dto.StoreActivityLogRequest true "Activity log"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.ValidationErrorResponse
// @Router /activity-logs [post]
// @Security BearerAuth
func (ctrl *ActivityLogController) Store(c *fiber.Ctx) error {
	var request dto.StoreActivityLogRequest

	rules := govalidator.MapData{
		"log_name":    []string{"required"},
		"description": []string{"required"},
		"event":       []string{"required"},
	}

	errs := utils.ValidateJSON(c, &request, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	activityLog := models.ActivityLog{
		LogName:        request.LogName,
		Description:    request.Description,
		SubjectID:      utils.Uint(request.SubjectID),
		SubjectType:    utils.String(request.SubjectType),
		Event:          request.Event,
		CauserID:       request.CauserID,
		CauserType:     request.CauserType,
		PrevProperties: &request.PrevProperties,
		Properties:     request.Properties,
		BatchUUID:      utils.String(request.BatchUUID),
		IPAddress:      utils.String(request.IPAddress),
		Country:        utils.String(request.Country),
		City:           utils.String(request.City),
		Region:         utils.String(request.Region),
		Postal:         utils.String(request.Postal),
		Geolocation:    utils.String(request.Geolocation),
		Timezone:       utils.String(request.Timezone),
		UserAgent:      utils.String(request.UserAgent),
		Referer:        utils.String(request.Referer),
	}

	if err := ctrl.repo.Store(&activityLog); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to store activity log")
	}

	return utils.SuccessResponse(c, "Activity log stored successfully", activityLog)
}
