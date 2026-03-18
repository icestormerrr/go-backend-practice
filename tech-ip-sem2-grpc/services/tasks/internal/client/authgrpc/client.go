package authgrpc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"tech-ip-sem2-grpc/proto/authpb"
	"tech-ip-sem2-grpc/shared/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrUnavailable  = errors.New("auth service unavailable")
	ErrUpstream     = errors.New("auth service error")
)

type Client struct {
	client  authpb.AuthServiceClient
	timeout time.Duration
	logger  *log.Logger
}

func New(addr string, timeout time.Duration, logger *log.Logger) (*Client, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return &Client{
		client:  authpb.NewAuthServiceClient(conn),
		timeout: timeout,
		logger:  logger,
	}, conn, nil
}

func (c *Client) Verify(ctx context.Context, authHeader string) error {
	token, ok := parseBearerToken(authHeader)
	if !ok {
		return ErrUnauthorized
	}

	verifyCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if requestID := middleware.GetRequestID(ctx); requestID != "" {
		verifyCtx = metadata.AppendToOutgoingContext(verifyCtx, "x-request-id", requestID)
	}

	c.logger.Printf("request_id=%s calling grpc verify timeout=%s", middleware.GetRequestID(ctx), c.timeout)
	_, err := c.client.Verify(verifyCtx, &authpb.VerifyRequest{Token: token})
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return ErrUpstream
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return ErrUnauthorized
	case codes.DeadlineExceeded, codes.Unavailable:
		return ErrUnavailable
	default:
		return ErrUpstream
	}
}

func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrUnavailable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusBadGateway
	}
}

func parseBearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(header[len(prefix):])
	if token == "" {
		return "", false
	}

	return token, true
}
