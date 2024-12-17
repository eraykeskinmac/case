// pkg/middleware/logger.go
package middleware

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var logger = log.New(os.Stdout, "\n🔍 [API] ", log.LstdFlags)

type RequestLog struct {
	ID           string      `json:"request_id"`
	Timestamp    time.Time   `json:"timestamp"`
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	Status       int         `json:"status"`
	Duration     float64     `json:"duration_ms"`
	IP           string      `json:"ip"`
	UserAgent    string      `json:"user_agent"`
	RequestBody  interface{} `json:"request_body,omitempty"`
	ResponseBody interface{} `json:"response_body,omitempty"`
	QueryParams  string      `json:"query_params,omitempty"`
	ErrorMessage string      `json:"error,omitempty"`
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)

		// Request body log
		var requestBody interface{}
		if len(c.Body()) > 0 {
			if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
				requestBody = string(c.Body())
			}
		}

		// Execute handler
		err := c.Next()

		// Response duration
		duration := time.Since(start).Milliseconds()

		// Log entry
		logEntry := RequestLog{
			ID:          requestID,
			Timestamp:   start,
			Method:      c.Method(),
			Path:        c.Path(),
			Status:      c.Response().StatusCode(),
			Duration:    float64(duration),
			IP:          c.IP(),
			UserAgent:   c.Get("User-Agent"),
			QueryParams: string(c.Request().URI().QueryString()),
			RequestBody: requestBody,
		}

		if err != nil {
			logEntry.ErrorMessage = err.Error()
		}

		// JSON olarak tek bir satırda yazdır
		entryJSON, _ := json.Marshal(logEntry)
		logger.Println(string(entryJSON))

		return err
	}
}
