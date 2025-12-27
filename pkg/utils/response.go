package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response structure - mirip dengan Laravel API Resource
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse structure
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

// Meta for pagination
type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
}

// SuccessResponse sends success response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse sends created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
	})
}

// PaginatedSuccessResponse sends paginated success response
func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, page, limit int, total int64) error {
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Meta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			LastPage:    lastPage,
		},
	})
}
