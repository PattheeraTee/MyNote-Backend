package gormRepository

import (
	"miw/entities"
	"gorm.io/gorm"
)

type GormTagRepository struct {
	db *gorm.DB
}

func NewGormTagRepository(db *gorm.DB) *GormTagRepository {
	return &GormTagRepository{db: db}
}

func (r *GormTagRepository) CreateTag(tag *entities.Tag) error {
	return r.db.Create(tag).Error
}

func (r *GormTagRepository) GetTagById(id uint) (*entities.Tag, error) {
	var tag entities.Tag
    if err := r.db.Preload("Notes", func(db *gorm.DB) *gorm.DB {
        return db.Select("notes.note_id")
    }).First(&tag, id).Error; err != nil {
        return nil, err
    }
    return &tag, nil
}

