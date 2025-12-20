package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	http_utils "github.com/icestormerrr/pz10-auth/internal/utils/http"
)

// RecoverMiddleware — middleware, который перехватывает паники,
// логирует их и возвращает 500 Internal Server Error.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC: %v\n%s", rec, debug.Stack())

				http_utils.WriteError(w, http.StatusInternalServerError, "internal server error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
