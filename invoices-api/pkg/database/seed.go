package database

import (
	"fmt"
	"invoices-api/internal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	var count int64
	db.Model(&models.Invoice{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	invoices := []models.Invoice{
		{
			ServiceName:   "DMP Service",
			InvoiceNumber: 1001,
			Date:          time.Date(2024, 3, 16, 0, 0, 0, 0, time.UTC),
			Amount:        1500.50,
			Status:        "Pending",
		},
		{
			ServiceName:   "SSP Service",
			InvoiceNumber: 1002,
			Date:          time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC),
			Amount:        2500.75,
			Status:        "Paid",
		},
		{
			ServiceName:   "DMP Service",
			InvoiceNumber: 1003,
			Date:          time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC),
			Amount:        750.25,
			Status:        "Unpaid",
		},
		{
			ServiceName:   "DDP Service",
			InvoiceNumber: 1004,
			Date:          time.Date(2024, 3, 16, 0, 0, 0, 0, time.UTC),
			Amount:        1500.50,
			Status:        "Pending",
		},
		{
			ServiceName:   "SSP Service",
			InvoiceNumber: 1005,
			Date:          time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC),
			Amount:        2500.75,
			Status:        "Paid",
		},
		{
			ServiceName:   "DMP Service",
			InvoiceNumber: 1006,
			Date:          time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC),
			Amount:        750.25,
			Status:        "Unpaid",
		},
		{
			ServiceName:   "SSP Service",
			InvoiceNumber: 1007,
			Date:          time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC),
			Amount:        750.25,
			Status:        "Unpaid",
		},
		{
			ServiceName:   "DSP Service",
			InvoiceNumber: 1008,
			Date:          time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC),
			Amount:        750.25,
			Status:        "Unpaid",
		},
	}

	log.Println("Seeding database...")
	if err := db.Create(&invoices).Error; err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	log.Printf("Successfully seeded %d invoices", len(invoices))
	return nil
}
