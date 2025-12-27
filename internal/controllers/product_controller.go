package controllers

import (
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ProductController handles product-related HTTP requests
type ProductController struct {
	repo *repositories.ProductRepository
}

// NewProductController creates new product controller instance
func NewProductController(repo *repositories.ProductRepository) *ProductController {
	return &ProductController{repo: repo}
}

// CreateProductRequest for validation
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

// UpdateProductRequest for validation
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

// Index retrieves all products with pagination - GET /api/products
func (ctrl *ProductController) Index(c *fiber.Ctx) error {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, total, err := ctrl.repo.FindAll(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve products")
	}

	return utils.PaginatedSuccessResponse(c, "Products retrieved successfully", products, page, limit, total)
}

// Show retrieves single product - GET /api/products/:id
func (ctrl *ProductController) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	product, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
	}

	return utils.SuccessResponse(c, "Product retrieved successfully", product)
}

// Store creates new product - POST /api/products
func (ctrl *ProductController) Store(c *fiber.Ctx) error {
	var req CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationError(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := ctrl.repo.Create(product); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product")
	}

	return utils.CreatedResponse(c, "Product created successfully", product)
}

// Update updates existing product - PUT /api/products/:id
func (ctrl *ProductController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationError(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	product, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Product not found")
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	if err := ctrl.repo.Update(product); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update product")
	}

	return utils.SuccessResponse(c, "Product updated successfully", product)
}

// Destroy deletes product - DELETE /api/products/:id
func (ctrl *ProductController) Destroy(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	if err := ctrl.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product")
	}

	return utils.SuccessResponse(c, "Product deleted successfully", nil)
}
