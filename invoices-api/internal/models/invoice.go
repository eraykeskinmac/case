package models

import "time"

type Invoice struct {
	ID            uint      `json:"id" gorm:"primaryKey" example:"1"`
	ServiceName   string    `json:"service_name" gorm:"not null" example:"DMP Service"`
	InvoiceNumber int       `json:"invoice_number" gorm:"unique" example:"1001"`
	Date          time.Time `json:"date" example:"2024-03-16T00:00:00Z"`
	Amount        float64   `json:"amount" example:"1500.50"`
	Status        string    `json:"status" example:"Pending"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
