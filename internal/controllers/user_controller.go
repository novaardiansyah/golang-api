package controllers

import (
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	repo *repositories.UserRepository
}

func NewUserController(repo *repositories.UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (ctrl *UserController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	total, err := ctrl.repo.Count()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count users")
	}

	users, err := ctrl.repo.FindAllPaginated(page, perPage)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users")
	}

	return utils.PaginatedSuccessResponse(c, "Users retrieved successfully", users, page, perPage, total, len(users))
}

func (ctrl *UserController) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
}
