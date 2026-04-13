package main

import (
	"log"
	"net/http"
	"os"

	"example.com/tech-ip-sem2-monitoring/internal/httpapi"
	"example.com/tech-ip-sem2-monitoring/internal/student"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	repo := student.NewRepo()
	handler := httpapi.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/students/", handler.GetStudentByID)
	mux.Handle("/metrics", promhttp.Handler())

	rootHandler := httpapi.MetricsMiddleware(mux)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("server started on %s", addr)
	if err := http.ListenAndServe(addr, rootHandler); err != nil {
		log.Fatal(err)
	}
}
