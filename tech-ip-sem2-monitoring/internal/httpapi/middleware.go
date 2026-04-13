package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"example.com/tech-ip-sem2-monitoring/internal/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start).Seconds()
		path := normalizePath(r.URL.Path)

		metrics.HttpRequestsTotal.WithLabelValues(r.Method, path).Inc()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

		if lrw.StatusCode() >= 400 {
			metrics.HttpErrorsTotal.WithLabelValues(
				r.Method,
				path,
				strconv.Itoa(lrw.StatusCode()),
			).Inc()
		}
	})
}

func normalizePath(path string) string {
	switch {
	case path == "/health":
		return "/health"
	case path == "/metrics":
		return "/metrics"
	case len(path) >= len("/students/") && path[:len("/students/")] == "/students/":
		return "/students/{id}"
	default:
		return path
	}
}
