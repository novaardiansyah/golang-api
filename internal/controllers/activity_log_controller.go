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
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
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
