package controllers

import (
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type GalleryController struct {
	repo *repositories.GalleryRepository
}

func NewGalleryController(repo *repositories.GalleryRepository) *GalleryController {
	return &GalleryController{repo: repo}
}

// Index godoc
// @Summary List galleries
// @Description Get a paginated list of galleries
// @Tags galleries
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]GallerySwagger}
// @Failure 400 {object} utils.Response
// @Router /galleries [get]
// @Security BearerAuth
func (ctrl *GalleryController) Index(c *fiber.Ctx) error {
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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to count galleries")
	}

	galleries, err := ctrl.repo.FindAllPaginated(page, perPage)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve galleries")
	}

	return utils.PaginatedSuccessResponse(c, "Galleries retrieved successfully", galleries, page, perPage, total, len(galleries))
}
