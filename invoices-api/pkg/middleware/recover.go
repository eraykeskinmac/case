package middleware

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()

				log.Printf("🔥 Panic recovered: %v\n%s", r, string(stack))

				errMsg := "Internal server error"
				if fiber.IsChild() { // Development ortamında
					errMsg = fmt.Sprintf("Panic: %v", r)
				}

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": errMsg,
					"error":   "A panic occurred in the server",
				})
			}
		}()

		return c.Next()
	}
}
