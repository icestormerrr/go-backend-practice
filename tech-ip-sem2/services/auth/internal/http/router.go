package httpapi

import (
	"log"
	"net/http"
)

func NewRouter(handler *Handler, _ *log.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/auth/login", handler.Login)
	mux.HandleFunc("/v1/auth/verify", handler.Verify)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}
