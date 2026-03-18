package grpcserver

import (
	"context"
	"log"

	"tech-ip-sem2-grpc/proto/authpb"
	"tech-ip-sem2-grpc/services/auth/internal/service"
	"tech-ip-sem2-grpc/shared/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	service *service.Service
	logger  *log.Logger
}

func Register(server *grpc.Server, service *service.Service, logger *log.Logger) {
	authpb.RegisterAuthServiceServer(server, &Server{
		service: service,
		logger:  logger,
	})
}

func (s *Server) Verify(ctx context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	subject, ok := s.service.VerifyToken(ctx, req.GetToken())
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	s.logger.Printf("request_id=%s verify_grpc subject=%s", middleware.GetRequestID(ctx), subject)
	return &authpb.VerifyResponse{
		Valid:   true,
		Subject: subject,
	}, nil
}

func RequestIDInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("x-request-id")
			if len(values) > 0 {
				requestID = values[0]
			}
		}

		if requestID == "" {
			requestID = "grpc-generated-request-id"
		}

		ctx = middleware.WithRequestID(ctx, requestID)
		logger.Printf("request_id=%s grpc_method=%s", requestID, info.FullMethod)
		return handler(ctx, req)
	}
}
