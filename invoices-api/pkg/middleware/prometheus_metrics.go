package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 3, 5, 7, 10},
		},
		[]string{"method", "path"},
	)
)

func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()

		err := c.Next()

		status := c.Response().StatusCode()
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(method, path, string(status)).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}

func PrometheusHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
		handler(c.Context())
		return nil
	}
}
