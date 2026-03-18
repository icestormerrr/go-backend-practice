package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/middleware"
)

type Handler struct {
	service *service.Service
	logger  *log.Logger
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type verifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *service.Service, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "username and password are required"})
		return
	}

	token, ok := h.service.Login(req.Username, req.Password)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid credentials"})
		return
	}

	h.logger.Printf("request_id=%s login subject=%s", middleware.GetRequestID(r.Context()), req.Username)
	writeJSON(w, http.StatusOK, loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	subject, ok := h.service.VerifyAuthorizationHeader(r.Header.Get("Authorization"))
	if !ok {
		writeJSON(w, http.StatusUnauthorized, verifyResponse{
			Valid: false,
			Error: "unauthorized",
		})
		return
	}

	h.logger.Printf("request_id=%s verify subject=%s", middleware.GetRequestID(r.Context()), subject)
	writeJSON(w, http.StatusOK, verifyResponse{
		Valid:   true,
		Subject: subject,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
