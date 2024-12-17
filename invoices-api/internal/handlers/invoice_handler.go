package handlers

import (
	"invoices-api/internal/models"
	"invoices-api/internal/repository"
	"invoices-api/pkg/validator"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler struct {
	repo      *repository.InvoiceRepository
	validator *validator.InvoiceValidator
}

func NewInvoiceHandler(repo *repository.InvoiceRepository, validator *validator.InvoiceValidator) *InvoiceHandler {
	return &InvoiceHandler{
		repo:      repo,
		validator: validator,
	}
}

func (h *InvoiceHandler) GetInvoices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "")
	sortDir := c.Query("sort_dir", "asc")

	var invoices []models.Invoice
	var total int64
	var err error

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	if search != "" {
		invoices, total, err = h.repo.Search(search, page, limit)
	} else {
		invoices, total, err = h.repo.GetAll(page, limit, sortBy, sortDir)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting invoices",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": invoices,
		"meta": fiber.Map{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
			"sort_by":     sortBy,
			"sort_dir":    sortDir,
		},
	})
}

func (h *InvoiceHandler) GetInvoiceByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID format",
			"details": "Please provide a valid invoice ID",
		})
	}

	invoice, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Invoice not found",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": invoice,
	})
}

func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	invoice := new(models.Invoice)

	if err := c.BodyParser(invoice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"details": "Please check your input data format",
		})
	}

	if errors := h.validator.ValidateInvoice(invoice); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	if err := h.repo.Create(invoice); err != nil {
		if strings.Contains(err.Error(), "invoice number") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Invoice number already exists",
				"details": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create invoice",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Invoice created successfully",
		"data":    invoice,
	})
}

func (h *InvoiceHandler) UpdateInvoice(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID format",
		})
	}

	existingInvoice, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Invoice not found",
		})
	}

	updates := make(map[string]interface{})
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errors := h.validator.ValidatePartialUpdate(updates); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	if status, exists := updates["status"]; exists && status != "" {
		existingInvoice.Status = status.(string)
	}
	if amount, exists := updates["amount"]; exists {
		existingInvoice.Amount = amount.(float64)
	}
	if serviceName, exists := updates["service_name"]; exists && serviceName != "" {
		existingInvoice.ServiceName = serviceName.(string)
	}
	if invoiceNumber, exists := updates["invoice_number"]; exists {
		existingInvoice.InvoiceNumber = int(invoiceNumber.(float64))
	}

	if err := h.repo.Update(&existingInvoice); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update invoice",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice updated successfully",
		"data":    existingInvoice,
	})
}

func (h *InvoiceHandler) DeleteInvoice(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID format",
			"details": "Please provide a valid invoice ID",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Invoice not found",
				"details": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete invoice",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice deleted successfully",
	})
}
