package core

type NoteCreatePayload struct {
	Title   string
	Content string
}

type NoteUpdatePayload struct {
	Title   *string
	Content *string
}

type NotesService interface {
	GetAll() ([]*Note, error)
	GetById(id int64) (Note, error)
	CreateNote(payload NoteCreatePayload) (int64, error)
	UpdateNote(id int64, payload NoteUpdatePayload) error
	Delete(id int64) error
}
