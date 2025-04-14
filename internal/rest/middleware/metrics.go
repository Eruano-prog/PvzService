package middleware

import (
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

type HTTPMetrics interface {
	RecordRequest(method, path, status string)
	RecordResponseTime(method, path string, duration float64)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Metrics(next http.Handler, metrics HTTPMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		statusText := http.StatusText(rw.statusCode)

		metrics.RecordRequest(r.Method, r.URL.Path, statusText)
		metrics.RecordResponseTime(r.Method, r.URL.Path, duration)
	})
}
