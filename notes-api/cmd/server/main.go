package main

import (
	"log"
	"net/http"

	router "github.com/icestormerrr/notes-api/internal/delivery/http"
	"github.com/icestormerrr/notes-api/internal/delivery/http/handlers"
	"github.com/icestormerrr/notes-api/internal/repos"
	"github.com/icestormerrr/notes-api/internal/services"
	"github.com/icestormerrr/notes-api/internal/utils/config"
)

// Package main Notes API server.
//
// @title           Notes API
// @version         1.0
// @description     Учебный REST API для заметок
// @contact.name    Backend Course
// @contact.email   yurkin.v.i@edu.mirea.ru
// @BasePath        /api/v1
func main() {
	cfg := config.Load()

	notesInmemoryRepo := repos.NewNoteRepoMem()
	notesService := services.NewNotesService(notesInmemoryRepo)
	notesHandler := handlers.NewNoteHandler(notesService)
	mux := router.Build(notesHandler)

	log.Println("listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
