package response

import "github.com/gofiber/fiber/v2"

// Meta holds pagination metadata.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success returns a 200 success response.
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}
	return c.JSON(Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Created returns a 201 created response.
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Created successfully"
	if len(message) > 0 {
		msg = message[0]
	}
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Paginated returns a paginated success response.
func Paginated(c *fiber.Ctx, data interface{}, meta Meta) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

// Error returns an error response with the given status code.
func Error(c *fiber.Ctx, statusCode int, message string, errors ...interface{}) error {
	resp := Response{
		Success: false,
		Message: message,
	}
	if len(errors) > 0 {
		resp.Errors = errors[0]
	}
	return c.Status(statusCode).JSON(resp)
}

// BadRequest returns a 400 error.
func BadRequest(c *fiber.Ctx, message string, errors ...interface{}) error {
	return Error(c, fiber.StatusBadRequest, message, errors...)
}

// Unauthorized returns a 401 error.
func Unauthorized(c *fiber.Ctx, message ...string) error {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusUnauthorized, msg)
}

// Forbidden returns a 403 error.
func Forbidden(c *fiber.Ctx, message ...string) error {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusForbidden, msg)
}

// NotFound returns a 404 error.
func NotFound(c *fiber.Ctx, message ...string) error {
	msg := "Resource not found"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusNotFound, msg)
}

// InternalError returns a 500 error.
func InternalError(c *fiber.Ctx, message ...string) error {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusInternalServerError, msg)
}
