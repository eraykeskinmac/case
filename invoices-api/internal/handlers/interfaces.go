package handlers

import "github.com/gofiber/fiber/v2"

type InvoiceHandler interface {
	GetInvoices(c *fiber.Ctx) error
	GetInvoiceByID(c *fiber.Ctx) error
	CreateInvoice(c *fiber.Ctx) error
	UpdateInvoice(c *fiber.Ctx) error
	DeleteInvoice(c *fiber.Ctx) error
}
