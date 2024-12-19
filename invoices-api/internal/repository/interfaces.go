package repository

import (
	"context"
	"invoices-api/internal/models"
)

type InvoiceRepository interface {
	GetAll(ctx context.Context, params QueryParams) ([]models.Invoice, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Invoice, error)
	Create(ctx context.Context, invoice *models.Invoice) error
	Update(ctx context.Context, invoice *models.Invoice) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, searchTerm string, params QueryParams) ([]models.Invoice, int64, error)
}

type QueryParams struct {
	Page    int
	Limit   int
	SortBy  string
	SortDir string
}

func NewQueryParams(page, limit int, sortBy, sortDir string) QueryParams {
	return QueryParams{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		SortDir: sortDir,
	}
}
