package handlers

import (
	"context"
	"fmt"
	"invoices-api/internal/models"
	"invoices-api/internal/repository"
	"invoices-api/pkg/middleware"
	"invoices-api/pkg/validator"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultPage    = 1
	defaultLimit   = 10
	maxLimit       = 100
	requestTimeout = 30 * time.Second
)

type RequestParams struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Search  string `json:"search"`
	SortBy  string `json:"sort_by"`
	SortDir string `json:"sort_dir"`
}

type invoiceHandler struct {
	repo      repository.InvoiceRepository
	validator *validator.InvoiceValidator
	cache     struct {
		sync.RWMutex
		data sync.Map
	}
}

func NewInvoiceHandler(repo repository.InvoiceRepository, validator *validator.InvoiceValidator) InvoiceHandler {
	return &invoiceHandler{
		repo:      repo,
		validator: validator,
	}
}

func (h *invoiceHandler) withTimeout(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Context(), requestTimeout)
}

func (h *invoiceHandler) GetInvoices(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()

	params := h.parseQueryParams(c)
	cacheKey := h.buildCacheKey("invoices", params)

	h.cache.RLock()
	cached, ok := h.cache.data.Load(cacheKey)
	h.cache.RUnlock()
	if ok {
		return c.JSON(cached)
	}

	var (
		invoices []models.Invoice
		total    int64
		err      error
	)

	queryParams := repository.NewQueryParams(params.Page, params.Limit, params.SortBy, params.SortDir)

	if params.Search != "" {
		invoices, total, err = h.repo.Search(ctx, params.Search, queryParams)
	} else {
		invoices, total, err = h.repo.GetAll(ctx, queryParams)
	}

	if err != nil {
		return err
	}

	response := fiber.Map{
		"data": invoices,
		"meta": h.buildMetadata(total, params),
	}

	h.cache.Lock()
	h.cache.data.Store(cacheKey, response)
	h.cache.Unlock()
	go h.scheduleInvalidateCache(cacheKey, 30*time.Second)

	return c.JSON(response)
}

func (h *invoiceHandler) GetInvoiceByID(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()

	id, err := h.parseID(c)
	if err != nil {
		return err
	}

	invoice, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": invoice,
	})
}

func (h *invoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()

	invoice := new(models.Invoice)
	if err := c.BodyParser(invoice); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	if errs := h.validator.ValidateInvoice(invoice); len(errs) > 0 {
		return middleware.NewBadRequestError("Validation failed", errs)
	}

	if err := h.repo.Create(ctx, invoice); err != nil {
		return err
	}

	h.invalidateListCache()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Invoice created successfully",
		"data":    invoice,
	})
}

func (h *invoiceHandler) UpdateInvoice(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()

	id, err := h.parseID(c)
	if err != nil {
		return err
	}

	invoice := new(models.Invoice)
	if err := c.BodyParser(invoice); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	if errs := h.validator.ValidateInvoice(invoice); len(errs) > 0 {
		return middleware.NewBadRequestError("Validation failed", errs)
	}

	invoice.ID = id
	if err := h.repo.Update(ctx, invoice); err != nil {
		return err
	}

	h.invalidateListCache()
	h.cache.data.Delete(h.buildCacheKey("invoice", id))

	return c.JSON(fiber.Map{
		"message": "Invoice updated successfully",
		"data":    invoice,
	})
}

func (h *invoiceHandler) DeleteInvoice(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()

	id, err := h.parseID(c)
	if err != nil {
		return err
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return err
	}

	h.invalidateListCache()
	h.cache.data.Delete(h.buildCacheKey("invoice", id))

	return c.JSON(fiber.Map{
		"message": "Invoice deleted successfully",
	})
}

func (h *invoiceHandler) parseQueryParams(c *fiber.Ctx) *RequestParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = defaultPage
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 || limit > maxLimit {
		limit = defaultLimit
	}

	sortDir := c.Query("sort_dir", "asc")
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc"
	}

	return &RequestParams{
		Page:    page,
		Limit:   limit,
		Search:  c.Query("search", ""),
		SortBy:  c.Query("sort_by", ""),
		SortDir: sortDir,
	}
}

func (h *invoiceHandler) parseID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return 0, middleware.NewBadRequestError("Invalid ID format")
	}
	return uint(id), nil
}

func (h *invoiceHandler) buildMetadata(total int64, params *RequestParams) fiber.Map {
	return fiber.Map{
		"total":       total,
		"page":        params.Page,
		"limit":       params.Limit,
		"total_pages": (total + int64(params.Limit) - 1) / int64(params.Limit),
		"sort_by":     params.SortBy,
		"sort_dir":    params.SortDir,
	}
}

func (h *invoiceHandler) buildCacheKey(prefix string, params interface{}) string {
	return fmt.Sprintf("%s:%v", prefix, params)
}

func (h *invoiceHandler) invalidateListCache() {
	h.cache.data.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok && strings.HasPrefix(k, "invoices:") {
			h.cache.data.Delete(key)
		}
		return true
	})
}

func (h *invoiceHandler) scheduleInvalidateCache(key interface{}, duration time.Duration) {
	time.Sleep(duration)
	h.cache.data.Delete(key)
}
