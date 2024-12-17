package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("error code: %d, message: %s", e.Code, e.Message)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	if e, ok := err.(*ErrorResponse); ok {
		code = e.Code
		message = e.Message
		return c.Status(code).JSON(e)
	}

	return c.Status(code).JSON(ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func NewBadRequestError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    fiber.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    fiber.StatusNotFound,
		Message: message,
	}
}

func NewInternalError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    fiber.StatusInternalServerError,
		Message: message,
	}
}
