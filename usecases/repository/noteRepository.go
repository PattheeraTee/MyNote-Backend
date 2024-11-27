package repository

import (
	"miw/entities"
)

type NoteRepository interface {
	CreateNote(note *entities.Note) error
	GetAllNoteByUserId(id uint) ([]entities.Note, error)
	GetNoteById(id uint) (*entities.Note, error)
	UpdateNote(note *entities.Note) error
	DeleteNoteById(id uint) error
	RestoreNoteById(id uint) error 
	AddTagToNote(noteID uint, tagID uint, userID uint) error
	RemoveTagFromNote(noteID uint, tagID uint, userID uint) error
	GetNoteByIdAndUser(noteID uint, userID uint) (*entities.Note, error)
}
