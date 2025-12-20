package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/icestormerrr/pz4-todo/internal/task"
	myMW "github.com/icestormerrr/pz4-todo/pkg/middleware"
)

func main() {
	repo := task.NewRepo("tasks.json")
	handler := task.NewHandler(repo)

	router := chi.NewRouter()
	router.Use(chimw.RequestID)
	router.Use(chimw.Recoverer)
	router.Use(myMW.Logger)
	router.Use(myMW.SimpleCORS)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/tasks", handler.Routes())
		})
	})

	addr := getAddr()
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
