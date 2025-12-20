package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/icestormerrr/notes-api/docs"
	"github.com/icestormerrr/notes-api/internal/delivery/http/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Build(notesHandler *handlers.NoteHandler) http.Handler {
	r := chi.NewRouter()

	// требование практической работы
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Get("/redoc", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/redoc.html")
	})
	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Get("/", notesHandler.GetAllNotes)
		r.Post("/", notesHandler.CreateNote)
		r.Get("/{id}", notesHandler.GetNoteById)
		r.Patch("/{id}", notesHandler.UpdateNote)
		r.Delete("/{id}", notesHandler.DeleteNote)
	})

	return r
}
