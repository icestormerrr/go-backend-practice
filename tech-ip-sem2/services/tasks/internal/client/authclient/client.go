package authclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"tech-ip-sem2/shared/middleware"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrUnavailable  = errors.New("auth service unavailable")
	ErrUpstream     = errors.New("auth service error")
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
	timeout    time.Duration
}

type verifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject"`
	Error   string `json:"error"`
}

func New(baseURL string, httpClient *http.Client, logger *log.Logger, timeout time.Duration) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
		logger:     logger,
		timeout:    timeout,
	}
}

func (c *Client) Verify(ctx context.Context, authHeader string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/auth/verify", nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpstream, err)
	}

	req.Header.Set("Authorization", authHeader)
	if requestID := middleware.GetRequestID(ctx); requestID != "" {
		req.Header.Set(middleware.RequestIDHeader, requestID)
	}

	c.logger.Printf("request_id=%s auth_verify url=%s timeout=%s", middleware.GetRequestID(ctx), req.URL.String(), c.timeout)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrUnavailable
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return ErrUnavailable
		}

		return ErrUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}

	var body verifyResponse
	_ = json.NewDecoder(resp.Body).Decode(&body)
	c.logger.Printf("request_id=%s auth_verify_failed status=%d error=%s", middleware.GetRequestID(ctx), resp.StatusCode, body.Error)

	if resp.StatusCode >= 500 {
		return ErrUnavailable
	}

	return ErrUpstream
}
