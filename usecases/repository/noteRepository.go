package repository

import(
	"miw/entities"
)

type NoteRepository interface {
	CreateNote(note *entities.Note) error
	UpdateNote(note *entities.Note) error
	GetAllNote(id uint) (*entities.Note, error)
	DeleteNoteById(id uint) error
}