package handlers

import "github.com/gofiber/fiber/v2"

func PanicTestHandler(c *fiber.Ctx) error {
	panic("Test panic!")
	return nil
}
