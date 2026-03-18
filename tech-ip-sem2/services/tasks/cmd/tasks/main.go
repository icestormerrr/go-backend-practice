package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	httpapi "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/service"
	"tech-ip-sem2/shared/httpx"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	port := getEnv("TASKS_PORT", "8082")
	authBaseURL := getEnv("AUTH_BASE_URL", "http://localhost:8081")
	authTimeout := getDurationMillis("AUTH_TIMEOUT_MS", 2500)

	logger := log.New(os.Stdout, "[tasks] ", log.LstdFlags|log.Lmicroseconds)
	taskService := service.New()
	authHTTPClient := httpx.NewClient(authTimeout)
	authClient := authclient.New(authBaseURL, authHTTPClient, logger, authTimeout)

	handler := httpapi.NewHandler(taskService, authClient, logger)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           withMiddleware(router, logger),
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Printf("tasks service started on :%s auth_base_url=%s auth_timeout=%s", port, authBaseURL, authTimeout)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("tasks service failed: %v", err)
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

func getDurationMillis(key string, fallbackMillis int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallbackMillis) * time.Millisecond
	}

	millis, err := strconv.Atoi(value)
	if err != nil || millis <= 0 {
		return time.Duration(fallbackMillis) * time.Millisecond
	}

	return time.Duration(millis) * time.Millisecond
}
