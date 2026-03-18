package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	"tech-ip-sem2/services/tasks/internal/service"
	"tech-ip-sem2/shared/middleware"
)

type Handler struct {
	service    *service.Service
	authClient *authclient.Client
	logger     *log.Logger
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *service.Service, authClient *authclient.Client, logger *log.Logger) *Handler {
	return &Handler{
		service:    service,
		authClient: authClient,
		logger:     logger,
	}
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	if !h.authorize(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.service.List())
	case http.MethodPost:
		h.createTask(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

func (h *Handler) TaskByID(w http.ResponseWriter, r *http.Request) {
	if !h.authorize(w, r) {
		return
	}

	id, ok := taskIDFromPath(r.URL.Path)
	if !ok {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := h.service.Get(id)
		if err != nil {
			h.writeTaskError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, task)
	case http.MethodPatch:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		if err := h.service.Delete(id); err != nil {
			h.writeTaskError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var input service.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	task, err := h.service.Create(input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	h.logger.Printf("request_id=%s task_created id=%s", middleware.GetRequestID(r.Context()), task.ID)
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var input service.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	task, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}

		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	h.logger.Printf("request_id=%s task_updated id=%s", middleware.GetRequestID(r.Context()), task.ID)
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "missing authorization header"})
		return false
	}

	err := h.authClient.Verify(r.Context(), authHeader)
	if err == nil {
		return true
	}

	switch {
	case errors.Is(err, authclient.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
	case errors.Is(err, authclient.ErrUnavailable):
		writeJSON(w, http.StatusServiceUnavailable, errorResponse{Error: "authorization service unavailable"})
	default:
		writeJSON(w, http.StatusBadGateway, errorResponse{Error: "authorization service error"})
	}

	return false
}

func (h *Handler) writeTaskError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
		return
	}

	writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
}

func taskIDFromPath(path string) (string, bool) {
	const prefix = "/v1/tasks/"
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}

	id := strings.TrimPrefix(path, prefix)
	if id == "" || strings.Contains(id, "/") {
		return "", false
	}

	return id, true
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
