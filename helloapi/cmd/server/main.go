package main

import (
	"log"
	"net/http"
	"os"

	"example.com/helloapi/internal/handlers"
	"example.com/helloapi/internal/middlewares"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", handlers.GetGreetingsHandler)
	mux.HandleFunc("/user", handlers.GetUserHandler)
	mux.HandleFunc("/health", handlers.GrtHealthCheckHandler)

	muxWithMiddlewares := middlewares.LoggingMiddleware(mux)
	addr := getAddr()

	log.Printf("Starting on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, muxWithMiddlewares))
}

func getAddr() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
