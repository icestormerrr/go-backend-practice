package core

type NotesRepository interface {
	GetAll() ([]*Note, error)
	GetById(id int64) (Note, error)
	Create(note Note) (int64, error)
	Update(note Note) error
	Delete(id int64) error
}
