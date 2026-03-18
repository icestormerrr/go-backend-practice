package main

import (
	"log"
	"net/http"
	"os"
	"time"

	httpapi "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	port := getEnv("AUTH_PORT", "8081")
	logger := log.New(os.Stdout, "[auth] ", log.LstdFlags|log.Lmicroseconds)

	authService := service.New()
	handler := httpapi.NewHandler(authService, logger)
	router := httpapi.NewRouter(handler, logger)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           withMiddleware(router, logger),
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Printf("auth service started on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("auth service failed: %v", err)
	}
}

func withMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return middleware.RequestID(middleware.Logging(logger)(next))
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
