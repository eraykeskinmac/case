package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var response ErrorResponse

	switch e := err.(type) {
	case *ErrorResponse:
		response = *e
	case *fiber.Error:
		response = ErrorResponse{
			Code:    e.Code,
			Message: e.Message,
		}
	default:
		response = ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return c.Status(response.Code).JSON(response)
}

func NewError(code int, message string, details ...interface{}) *ErrorResponse {
	err := &ErrorResponse{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

func NewBadRequestError(message string, details ...interface{}) *ErrorResponse {
	return NewError(fiber.StatusBadRequest, message, details...)
}

func NewNotFoundError(message string, details ...interface{}) *ErrorResponse {
	return NewError(fiber.StatusNotFound, message, details...)
}

func NewInternalError(message string, details ...interface{}) *ErrorResponse {
	return NewError(fiber.StatusInternalServerError, message, details...)
}
