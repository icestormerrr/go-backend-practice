package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func LoggingMiddleware(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := time.Now().UnixNano()
		lrw := NewLoggingResponseWriter(w)

		log.Info("incoming request",
			zap.Int64("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)

		lrw.Header().Set("X-Request-ID", strconv.FormatInt(requestID, 10))
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Info("request completed",
			zap.Int64("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status_code", lrw.StatusCode()),
			zap.Duration("duration", duration),
		)
	})
}
