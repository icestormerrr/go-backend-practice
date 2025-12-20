package middleware

import (
	"net/http"

	http_utils "github.com/icestormerrr/pz10-auth/internal/utils/http"
)

func AuthZRoles(allowed ...string) func(http.Handler) http.Handler {
	set := map[string]struct{}{}
	for _, a := range allowed {
		set[a] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, _ := r.Context().Value(ctxClaimsKey).(map[string]any)
			role, _ := claims["role"].(string)
			if _, ok := set[role]; !ok {
				http_utils.WriteError(w, http.StatusForbidden, "forbidden", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
