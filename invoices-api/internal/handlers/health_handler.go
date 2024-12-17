package handlers

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

type HealthResponse struct {
	Status     string            `json:"status"`
	Components map[string]Status `json:"components"`
	Timestamp  string            `json:"timestamp"`
}

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	health := HealthResponse{
		Status:     "healthy",
		Components: make(map[string]Status),
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	if err := h.checkDB(); err != nil {
		health.Status = "unhealthy"
		health.Components["database"] = Status{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	} else {
		health.Components["database"] = Status{
			Status: "healthy",
		}
	}

	memStatus := h.checkMemory()
	health.Components["memory"] = memStatus
	if memStatus.Status == "unhealthy" {
		health.Status = "unhealthy"
	}

	return c.JSON(health)
}

func (h *HealthHandler) checkDB() error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (h *HealthHandler) checkMemory() Status {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// high memory
	if m.Alloc > 1024*1024*1024 {
		return Status{
			Status:  "unhealthy",
			Message: "High memory usage",
		}
	}

	return Status{
		Status: "healthy",
	}
}
