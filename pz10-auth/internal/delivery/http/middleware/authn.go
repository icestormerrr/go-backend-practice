package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/icestormerrr/pz10-auth/internal/core"
	http_utils "github.com/icestormerrr/pz10-auth/internal/utils/http"
)

const ctxClaimsKey string = "claims"

func AuthN(v core.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				http_utils.WriteError(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}
			raw := strings.TrimPrefix(h, "Bearer ")
			claims, err := v.Parse(raw)
			if err != nil {
				http_utils.WriteError(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}
			ctx := context.WithValue(r.Context(), ctxClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
