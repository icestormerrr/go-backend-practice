package repos

import (
	"errors"
	"sync"

	"github.com/icestormerrr/notes-api/internal/core"
)

type NotesInmemoryRepo struct {
	mu    sync.RWMutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NotesInmemoryRepo {
	return &NotesInmemoryRepo{notes: make(map[int64]*core.Note)}
}

func (r *NotesInmemoryRepo) GetAll() ([]*core.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*core.Note
	for _, note := range r.notes {
		result = append(result, note)
	}

	return result, nil
}

func (r *NotesInmemoryRepo) GetById(id int64) (core.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result, ok := r.notes[id]
	if !ok {
		return core.Note{}, errors.New("note not found")
	}

	return *result, nil
}

func (r *NotesInmemoryRepo) Create(noteToCreate core.Note) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.next++
	noteToCreate.ID = r.next
	r.notes[noteToCreate.ID] = &noteToCreate

	return noteToCreate.ID, nil
}

func (r *NotesInmemoryRepo) Update(noteToUpdate core.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.notes[noteToUpdate.ID] = &noteToUpdate

	return nil
}

func (r *NotesInmemoryRepo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.notes[id]
	if !ok {
		return errors.New("note not found")
	}
	delete(r.notes, id)
	return nil
}
