package repository

import (
	"miw/entities"
)

type NoteRepository interface {
	CreateNote(note *entities.Note) error
	UpdateNote(note *entities.Note) error
	GetAllNoteByUserId(id uint) ([]entities.Note, error)
	GetNoteById(id uint) (*entities.Note, error)
	DeleteNoteById(id uint) error
	AddTagToNote(noteID uint, tagID uint) error // รองรับ tagID ทีละตัว
	// CreateTag(tag *entities.Tag) error
	// AddTagToNote(noteId uint, tagId uint) error
	// GetTagById(id uint) (*entities.Tag, error)
}
