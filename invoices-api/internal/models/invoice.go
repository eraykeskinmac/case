package models

import "time"

type Invoice struct {
	ID            uint      `json:"id" gorm:"primaryKey;column:id"`
	ServiceName   string    `json:"service_name" gorm:"column:service_name;not null"`
	InvoiceNumber int       `json:"invoice_number" gorm:"column:invoice_number;unique"`
	Date          time.Time `json:"date" gorm:"column:date"`
	Amount        float64   `json:"amount" gorm:"column:amount"`
	Status        string    `json:"status" gorm:"column:status"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Invoice) TableName() string {
	return "invoices"
}

type InvoiceResponse struct {
	Message string  `json:"message" example:"Operation successful"`
	Data    Invoice `json:"data"`
}

type InvoiceListResponse struct {
	Data []Invoice `json:"data"`
	Meta MetaData  `json:"meta"`
}

type MetaData struct {
	Total      int64  `json:"total" example:"100"`
	Page       int    `json:"page" example:"1"`
	Limit      int    `json:"limit" example:"10"`
	TotalPages int64  `json:"total_pages" example:"10"`
	SortBy     string `json:"sort_by,omitempty" example:"created_at"`
	SortDir    string `json:"sort_dir,omitempty" example:"desc"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error occurred"`
	Error   string `json:"error,omitempty" example:"Invalid input"`
}
