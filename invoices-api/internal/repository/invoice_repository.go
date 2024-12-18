package repository

import (
	"errors"
	"fmt"
	"invoices-api/internal/models"
	"strconv"

	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

type DuplicateInvoiceError struct {
	InvoiceNumber int
}

func (e *DuplicateInvoiceError) Error() string {
	return fmt.Sprintf("Invoice number %d already exists", e.InvoiceNumber)
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) GetAll(page int, limit int, sortBy string, sortDir string) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&models.Invoice{})

	query.Count(&total)

	if sortBy != "" {
		if sortDir == "desc" {
			query = query.Order(sortBy + " DESC")
		} else {
			query = query.Order(sortBy + " ASC")
		}
	} else {
		query = query.Order("id ASC")
	}

	result := query.Offset(offset).Limit(limit).Find(&invoices)

	return invoices, total, result.Error
}

func (r *InvoiceRepository) GetByID(id uint) (models.Invoice, error) {
	var invoice models.Invoice
	result := r.db.First(&invoice, id)
	return invoice, result.Error
}

func (r *InvoiceRepository) Create(invoice *models.Invoice) error {
	var existingInvoice models.Invoice
	result := r.db.Where("invoice_number = ?", invoice.InvoiceNumber).First(&existingInvoice)

	if result.Error == nil {
		return fmt.Errorf("invoice number %d already exists", invoice.InvoiceNumber)
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return r.db.Create(invoice).Error
}

func (r *InvoiceRepository) Update(invoice *models.Invoice) error {
	var existingInvoice models.Invoice
	result := r.db.First(&existingInvoice, invoice.ID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("invoice not found with id: %d", invoice.ID)
		}
		return result.Error
	}

	if invoice.InvoiceNumber != existingInvoice.InvoiceNumber {
		var exists bool
		r.db.Model(&models.Invoice{}).
			Select("count(*) > 0").
			Where("invoice_number = ? AND id != ?", invoice.InvoiceNumber, invoice.ID).
			Find(&exists)

		if exists {
			return &DuplicateInvoiceError{InvoiceNumber: invoice.InvoiceNumber}
		}
	}

	updates := make(map[string]interface{})

	if invoice.ServiceName != "" {
		updates["service_name"] = invoice.ServiceName
	}
	if invoice.InvoiceNumber != 0 {
		updates["invoice_number"] = invoice.InvoiceNumber
	}
	if !invoice.Date.IsZero() {
		updates["date"] = invoice.Date
	}
	if invoice.Amount > 0 {
		updates["amount"] = invoice.Amount
	}
	if invoice.Status != "" {
		validStatuses := map[string]bool{
			"Paid":    true,
			"Pending": true,
			"Unpaid":  true,
		}
		if !validStatuses[invoice.Status] {
			return fmt.Errorf("invalid status: %s", invoice.Status)
		}
		updates["status"] = invoice.Status
	}

	return r.db.Model(&existingInvoice).Updates(updates).Error
}

func (r *InvoiceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Invoice{}, id).Error
}

func (r *InvoiceRepository) Search(searchTerm string, page int, limit int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	offset := (page - 1) * limit

	invoiceNum, err := strconv.Atoi(searchTerm)
	var query *gorm.DB

	if err == nil {
		query = r.db.Model(&models.Invoice{}).Where("invoice_number = ?", invoiceNum)
	} else {
		query = r.db.Model(&models.Invoice{}).Where("service_name ILIKE ?", "%"+searchTerm+"%")
	}

	query.Count(&total)

	result := query.Offset(offset).Limit(limit).Find(&invoices)

	return invoices, total, result.Error
}
