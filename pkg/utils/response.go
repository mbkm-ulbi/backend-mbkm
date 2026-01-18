package utils

import (
	"github.com/gofiber/fiber/v2"
)

// StandardResponse represents a standard API response
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
	Page  int         `json:"page,omitempty"`
	Limit int         `json:"per_page,omitempty"`
}

// SuccessResponse returns a success response
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse returns an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors interface{}) error {
	return c.Status(statusCode).JSON(StandardResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// PaginatedSuccessResponse returns a paginated success response
func PaginatedSuccessResponse(c *fiber.Ctx, data interface{}, count int64) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Data:  data,
		Count: count,
	})
}

// ValidationError returns a validation error response
func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", errors)
}

// UnauthorizedError returns an unauthorized error response
func UnauthorizedError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return ErrorResponse(c, fiber.StatusUnauthorized, message, nil)
}

// ForbiddenError returns a forbidden error response
func ForbiddenError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "403 Forbidden"
	}
	return ErrorResponse(c, fiber.StatusForbidden, message, nil)
}

// NotFoundError returns a not found error response
func NotFoundError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return ErrorResponse(c, fiber.StatusNotFound, message, nil)
}

// InternalServerError returns an internal server error response
func InternalServerError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return ErrorResponse(c, fiber.StatusInternalServerError, message, nil)
}
