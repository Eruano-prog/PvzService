package monitoring

import (
	"AvitoPvz/internal/rest/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusHTTPMetrics struct {
	requestsTotal *prometheus.CounterVec
	responseTime  *prometheus.HistogramVec
}

func NewPrometheusHTTPMetrics() middleware.HTTPMetrics {
	return &PrometheusHTTPMetrics{
		requestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		}, []string{"method", "path", "status"}),

		responseTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		}, []string{"method", "path"}),
	}
}

func (m *PrometheusHTTPMetrics) RecordRequest(method, path, status string) {
	m.requestsTotal.WithLabelValues(method, path, status).Inc()
}

func (m *PrometheusHTTPMetrics) RecordResponseTime(method, path string, duration float64) {
	m.responseTime.WithLabelValues(method, path).Observe(duration)
}
