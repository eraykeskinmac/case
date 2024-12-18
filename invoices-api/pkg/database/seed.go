package database

import (
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
			ServiceName:   "SSP Service",
			InvoiceNumber: 1004,
			Date:          time.Date(2024, 3, 19, 0, 0, 0, 0, time.UTC),
			Amount:        3500.00,
			Status:        "Pending",
		},
		{
			ServiceName:   "DMP Service",
			InvoiceNumber: 1005,
			Date:          time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
			Amount:        1200.00,
			Status:        "Paid",
		},
	}

	log.Println("Seeding database...")
	result := db.Create(&invoices)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Successfully seeded %d invoices", len(invoices))
	return nil
}
