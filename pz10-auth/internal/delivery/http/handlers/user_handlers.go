package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/icestormerrr/pz10-auth/internal/core"
	http_utils "github.com/icestormerrr/pz10-auth/internal/utils/http"
)

type UserHandler struct {
	userService core.UserService
}

func NewUserHandler(userSvc core.UserService) *UserHandler {
	return &UserHandler{userService: userSvc}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(core.CtxClaimsKey).(map[string]any)
	if !ok {
		http_utils.WriteError(w, http.StatusUnauthorized, "invalid_claims", nil)
		return
	}
	http_utils.WriteJSON(w, map[string]any{
		"id": claims["sub"], "email": claims["email"], "role": claims["role"],
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid_user_id", nil)
		return
	}

	claims, ok := r.Context().Value(core.CtxClaimsKey).(map[string]any)
	if !ok {
		http_utils.WriteError(w, http.StatusUnauthorized, "invalid_claims", nil)
		return
	}

	if userID != int64(claims["sub"].(float64)) && claims["role"] != "admin" {
		http_utils.WriteError(w, http.StatusForbidden, "forbidden", nil)
		return
	}

	user, err := h.userService.GetById(userID)
	if err != nil {
		http_utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	http_utils.WriteJSON(w, user)
}

func (h *UserHandler) GetAdminStats(w http.ResponseWriter, r *http.Request) {
	http_utils.WriteJSON(w, map[string]any{"users": 2, "version": "1.0"})
}
