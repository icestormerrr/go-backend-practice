package services

import (
	"errors"
	"time"

	"github.com/icestormerrr/notes-api/internal/core"
)

type NotesService struct {
	notesRepo core.NotesRepository
}

func NewNotesService(notesRepo core.NotesRepository) *NotesService {
	return &NotesService{notesRepo: notesRepo}
}

func (s *NotesService) GetAll() ([]*core.Note, error) {
	return s.notesRepo.GetAll()
}

func (s *NotesService) GetById(id int64) (core.Note, error) {
	return s.notesRepo.GetById(id)
}

func (s *NotesService) CreateNote(payload core.NoteCreatePayload) (int64, error) {
	if payload.Title == "" {
		return 0, errors.New("title cannot be empty")
	}

	now := time.Now()

	n := core.Note{
		Title:     payload.Title,
		Content:   payload.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.notesRepo.Create(n)
}

func (s *NotesService) UpdateNote(id int64, payload core.NoteUpdatePayload) error {
	note, err := s.notesRepo.GetById(id)
	if err != nil {
		return errors.New("note not found")
	}

	if payload.Title != nil && *payload.Title == "" {
		return errors.New("title cannot be empty")
	}

	if payload.Title != nil {
		note.Title = *payload.Title
	}

	if payload.Content != nil {
		note.Content = *payload.Content
	}

	note.UpdatedAt = time.Now()

	return s.notesRepo.Update(note)
}

func (s *NotesService) Delete(id int64) error {
	return s.notesRepo.Delete(id)
}
