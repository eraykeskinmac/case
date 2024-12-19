package repository

import (
	"context"
	"fmt"
	"invoices-api/internal/models"
	"invoices-api/pkg/middleware"
	"sync"
	"time"

	"gorm.io/gorm"
)

const (
	defaultTimeout = 10 * time.Second
	maxSearchLen   = 100
)

type invoiceRepository struct {
	db    *gorm.DB
	cache sync.Map
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{
		db: db,
	}
}

func (r *invoiceRepository) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
}

func (r *invoiceRepository) GetAll(ctx context.Context, params QueryParams) ([]models.Invoice, int64, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var invoices []models.Invoice
	var total int64

	queryCount := r.db.WithContext(ctx).Model(&models.Invoice{})
	queryFetch := r.db.WithContext(ctx).Model(&models.Invoice{})

	if err := queryCount.Count(&total).Error; err != nil {
		return nil, 0, middleware.NewInternalError("Failed to count invoices")
	}

	if params.SortBy != "" {
		order := fmt.Sprintf("%s %s", params.SortBy, params.SortDir)
		queryFetch = queryFetch.Order(order)
	}

	offset := (params.Page - 1) * params.Limit
	queryFetch = queryFetch.Offset(offset).Limit(params.Limit)

	if err := queryFetch.Select("id, service_name, invoice_number, date, amount, status, created_at, updated_at").
		Find(&invoices).Error; err != nil {
		return nil, 0, middleware.NewInternalError("Failed to fetch invoices")
	}

	return invoices, total, nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uint) (*models.Invoice, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	if cached, ok := r.cache.Load(id); ok {
		if invoice, ok := cached.(*models.Invoice); ok {
			return invoice, nil
		}
	}

	var invoice models.Invoice
	if err := r.db.WithContext(ctx).First(&invoice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, middleware.NewNotFoundError("Invoice not found")
		}
		return nil, middleware.NewInternalError("Failed to fetch invoice")
	}

	r.cache.Store(id, &invoice)

	return &invoice, nil
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *models.Invoice) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.checkDuplicateInvoiceNumber(ctx, tx, invoice.InvoiceNumber); err != nil {
			return err
		}

		if err := tx.Create(invoice).Error; err != nil {
			return middleware.NewInternalError("Failed to create invoice")
		}

		return nil
	})
}

func (r *invoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.Invoice
		if err := tx.First(&existing, invoice.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return middleware.NewNotFoundError("Invoice not found")
			}
			return middleware.NewInternalError("Failed to fetch invoice")
		}

		if invoice.InvoiceNumber != existing.InvoiceNumber {
			if err := r.checkDuplicateInvoiceNumber(ctx, tx, invoice.InvoiceNumber, invoice.ID); err != nil {
				return err
			}
		}

		if err := tx.Save(invoice).Error; err != nil {
			return middleware.NewInternalError("Failed to update invoice")
		}

		r.cache.Delete(invoice.ID)

		return nil
	})
}

func (r *invoiceRepository) Delete(ctx context.Context, id uint) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&models.Invoice{}, id)
	if result.Error != nil {
		return middleware.NewInternalError("Failed to delete invoice")
	}

	if result.RowsAffected == 0 {
		return middleware.NewNotFoundError("Invoice not found")
	}

	r.cache.Delete(id)

	return nil
}

func (r *invoiceRepository) Search(ctx context.Context, searchTerm string, params QueryParams) ([]models.Invoice, int64, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var invoices []models.Invoice
	var total int64

	if len(searchTerm) > maxSearchLen {
		searchTerm = searchTerm[:maxSearchLen]
	}

	query := r.db.WithContext(ctx).Model(&models.Invoice{})

	if searchTerm != "" {
		query = query.Where("service_name ILIKE ?", fmt.Sprintf("%%%s%%", searchTerm))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, middleware.NewInternalError("Failed to count search results")
	}

	if params.SortBy != "" {
		order := fmt.Sprintf("%s %s", params.SortBy, params.SortDir)
		query = query.Order(order)
	}

	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).
		Limit(params.Limit).
		Select("id, service_name, invoice_number, date, amount, status, created_at, updated_at").
		Find(&invoices).Error; err != nil {
		return nil, 0, middleware.NewInternalError("Failed to fetch search results")
	}

	return invoices, total, nil
}

func (r *invoiceRepository) checkDuplicateInvoiceNumber(ctx context.Context, tx *gorm.DB, invoiceNumber int, excludeID ...uint) error {
	query := tx.WithContext(ctx).Model(&models.Invoice{}).Where("invoice_number = ?", invoiceNumber)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return middleware.NewInternalError("Failed to check duplicate invoice number")
	}

	if count > 0 {
		return middleware.NewBadRequestError(fmt.Sprintf("Invoice number %d already exists", invoiceNumber))
	}

	return nil
}
