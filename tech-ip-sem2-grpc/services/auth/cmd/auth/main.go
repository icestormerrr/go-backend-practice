package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcserver "tech-ip-sem2-grpc/services/auth/internal/grpc"
	httpapi "tech-ip-sem2-grpc/services/auth/internal/http"
	"tech-ip-sem2-grpc/services/auth/internal/service"
	"tech-ip-sem2-grpc/shared/middleware"

	"google.golang.org/grpc"
)

func main() {
	httpPort := getEnv("AUTH_HTTP_PORT", "8081")
	grpcPort := getEnv("AUTH_GRPC_PORT", "50051")
	logger := log.New(os.Stdout, "[auth] ", log.LstdFlags|log.Lmicroseconds)

	authService := service.New()
	httpHandler := httpapi.NewHandler(authService, logger)
	httpRouter := httpapi.NewRouter(httpHandler)

	httpServer := &http.Server{
		Addr:              ":" + httpPort,
		Handler:           middleware.RequestID(middleware.Logging(logger)(httpRouter)),
		ReadHeaderTimeout: 3 * time.Second,
	}

	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Fatalf("failed to listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcserver.RequestIDInterceptor(logger)),
	)
	grpcserver.Register(grpcServer, authService, logger)

	go func() {
		logger.Printf("auth HTTP service started on :%s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("auth HTTP service failed: %v", err)
		}
	}()

	go func() {
		logger.Printf("auth gRPC service started on :%s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatalf("auth gRPC service failed: %v", err)
		}
	}()

	waitForShutdown(logger, httpServer, grpcServer)
}

func waitForShutdown(logger *log.Logger, httpServer *http.Server, grpcServer *grpc.Server) {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	<-stopSignal

	logger.Println("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Printf("http shutdown error: %v", err)
	}

	grpcServer.GracefulStop()
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
