package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	grpcclient "tech-ip-sem2-grpc/services/tasks/internal/client/authgrpc"
	httpapi "tech-ip-sem2-grpc/services/tasks/internal/http"
	"tech-ip-sem2-grpc/services/tasks/internal/service"
	"tech-ip-sem2-grpc/shared/middleware"
)

func main() {
	port := getEnv("TASKS_PORT", "8082")
	authGRPCAddr := getEnv("AUTH_GRPC_ADDR", "localhost:50051")
	authTimeout := getDurationMillis("AUTH_GRPC_TIMEOUT_MS", 1500)
	logger := log.New(os.Stdout, "[tasks] ", log.LstdFlags|log.Lmicroseconds)

	authClient, conn, err := grpcclient.New(authGRPCAddr, authTimeout, logger)
	if err != nil {
		logger.Fatalf("failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	taskService := service.New()
	handler := httpapi.NewHandler(taskService, authClient, logger)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           middleware.RequestID(middleware.Logging(logger)(router)),
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		logger.Printf("tasks HTTP service started on :%s auth_grpc_addr=%s auth_timeout=%s", port, authGRPCAddr, authTimeout)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("tasks HTTP service failed: %v", err)
		}
	}()

	waitForShutdown(logger, server)
}

func waitForShutdown(logger *log.Logger, server *http.Server) {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	<-stopSignal

	logger.Println("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("http shutdown error: %v", err)
	}
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
