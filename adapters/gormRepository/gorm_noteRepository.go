package gormRepository

import (
	"gorm.io/gorm"
	"miw/entities"
)

type GormNoteRepository struct {
	db *gorm.DB
}

func NewGormNoteRepository(db *gorm.DB) *GormNoteRepository {
	return &GormNoteRepository{db: db}
}

func (r *GormNoteRepository) CreateNote(note *entities.Note) error {
	return r.db.Create(note).Error
}

func (r *GormNoteRepository) UpdateNote(note *entities.Note) error {
	return r.db.Save(note).Error
}

func (r *GormNoteRepository) GetAllNote(id uint) (*entities.Note, error) {
	var note entities.Note
	if err := r.db.First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GormNoteRepository) DeleteNoteById(id uint) error {
	return r.db.Delete(&entities.Note{}, id).Error
}
