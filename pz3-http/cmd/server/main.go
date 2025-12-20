package main

import (
	"log"
	"net/http"
	"os"

	"github.com/icestormerrr/pz3-http/internal/api"
	"github.com/icestormerrr/pz3-http/internal/storage"
)

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("PATCH /tasks/", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/", h.DeleteTask)
	mux.HandleFunc("GET /tasks/", h.GetTask)

	handler := api.WithCORS(api.WithLogging(mux))
	addr := getAddr()
	log.Println("listening on", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
