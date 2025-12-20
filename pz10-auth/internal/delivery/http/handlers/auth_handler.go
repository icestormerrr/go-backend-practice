package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/icestormerrr/pz10-auth/internal/core"
	http_utils "github.com/icestormerrr/pz10-auth/internal/utils/http"
)

type AuthHandler struct {
	authService core.AuthService
}

func NewAuthHandler(authSvc core.AuthService) *AuthHandler {
	return &AuthHandler{authService: authSvc}
}

// TODO: Простейший rate limit для /login (например, 5 попыток за 5 минут по IP).
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid_credentials", nil)
		return
	}

	accessToken, refreshToken, userID, err := h.authService.Login(in.Email, in.Password)
	if err != nil {
		http_utils.WriteError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	http_utils.WriteJSON(w, map[string]any{"token": accessToken, "userID": userID})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http_utils.WriteError(w, http.StatusUnauthorized, "missing_refresh_token", nil)
		return
	}

	newAccess, newRefresh, err := h.authService.RefreshTokens(cookie.Value)
	if err != nil {
		http_utils.WriteError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	newCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    newRefresh,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}
	http.SetCookie(w, &newCookie)

	http_utils.WriteJSON(w, map[string]any{"token": newAccess})
}
