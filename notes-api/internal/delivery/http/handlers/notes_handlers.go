package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/icestormerrr/notes-api/internal/core"
	http_utils "github.com/icestormerrr/notes-api/internal/utils/http"
)

type NoteHandler struct {
	service core.NotesService
}

func NewNoteHandler(service core.NotesService) *NoteHandler {
	return &NoteHandler{service: service}
}

type NoteCreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteUpdateRequest struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

// GetAllNotes godoc
// @Summary      Список заметок
// @Description  Возвращает список заметок с пагинацией и фильтром по заголовку
// @Tags         notes
// @Success      200    {array}  core.Note
// @Failure      500    {object}  map[string]string
// @Router       /notes [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetAll()
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "failed to get notes", err.Error())
		return
	}

	http_utils.WriteJSON(w, notes)
}

// GetNoteById godoc
// @Summary      Получить заметку
// @Tags         notes
// @Param        id   path   int  true  "ID"
// @Success      200  {object}  core.Note
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [get]
func (h *NoteHandler) GetNoteById(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	note, err := h.service.GetById(id)
	if err != nil {
		http_utils.WriteError(w, http.StatusNotFound, "note not found", nil)
		return
	}

	http_utils.WriteJSON(w, note)
}

// CreateNote godoc
// @Summary      Создать заметку
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        input  body     NoteCreateRequest  true  "Данные новой заметки"
// @Success      201    {object} core.Note
// @Failure      400    {object} map[string]string
// @Failure      500    {object} map[string]string
// @Router       /notes [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req NoteCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	payload := core.NoteCreatePayload{
		Title:   req.Title,
		Content: req.Content,
	}

	id, err := h.service.CreateNote(payload)
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	http_utils.WriteJSON(w, map[string]int64{"id": id})
}

// UpdateNote godoc
// @Summary      Обновить заметку (частично)
// @Tags         notes
// @Accept       json
// @Param        id     path   int        true  "Идентификатор заметки"
// @Param        input  body   NoteUpdateRequest true  "Поля для обновления"
// @Success      200    {object}  core.Note
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /notes/{id} [patch]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	var req NoteUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	payload := core.NoteUpdatePayload{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.service.UpdateNote(id, payload); err != nil {
		if err.Error() == "note not found" {
			http_utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		http_utils.WriteError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteNote godoc
// @Summary      Удалить заметку
// @Tags         notes
// @Param        id  path  int  true  "Bltynbabrfnjh pfvtnrb"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	if err := h.service.Delete(id); err != nil {
		http_utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}
