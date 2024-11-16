package gormRepository

import (
	"miw/entities"

	"gorm.io/gorm"
)

type GormNoteRepository struct {
	db *gorm.DB
}

func NewGormNoteRepository(db *gorm.DB) *GormNoteRepository {
	return &GormNoteRepository{db: db}
}

func (r *GormNoteRepository) CreateNote(note *entities.Note) error {
	// บันทึก Note ลงในฐานข้อมูล
	if err := r.db.Create(note).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormNoteRepository) UpdateNote(note *entities.Note) error {
	return r.db.Save(note).Error
}

func (r *GormNoteRepository) GetAllNoteByUserId(userID uint) ([]entities.Note, error) {
	var notes []entities.Note
    if err := r.db.Where("user_id = ?", userID).
        Preload("Tags", func(db *gorm.DB) *gorm.DB {
            return db.Select("tag_id, tag_name") // ไม่ดึง Notes ใน Tags
        }).
        Preload("Reminders").
        Preload("Event").
        Find(&notes).Error; err != nil {
        return nil, err
    }
	return notes, nil
}

// func (r *GormNoteRepository) GetNoteById(id uint) (*entities.Note, error) {
// 	var note entities.Note
// 	if err := r.db.Where("user_id = ?", id).
// 	Preload("Tags", func(db *gorm.DB) *gorm.DB {
// 		return db.Select("tag_id, tag_name") // ไม่ดึง Notes ใน Tags
// 	}).
// 	Preload("Reminders").
// 	Preload("Event").
// 	First(&note, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &note, nil
// }
func (r *GormNoteRepository) GetNoteById(id uint) (*entities.Note, error) {
	var note entities.Note
	if err := r.db.Where("note_id = ?", id).
	Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("tag_id, tag_name") // ไม่ดึง Notes ใน Tags
	}).
	Preload("Reminders").
	Preload("Event").
	First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GormNoteRepository) DeleteNoteById(id uint) error {
	return r.db.Delete(&entities.Note{}, id).Error
}


func (r *GormNoteRepository) AddTagToNote(noteID uint, tagID uint) error {
	var note entities.Note
	if err := r.db.Preload("Tags").First(&note, noteID).Error; err != nil {
		return err
	}
	var tag entities.Tag
	if err := r.db.First(&tag, tagID).Error; err != nil {
		return err
	}
	return r.db.Model(&note).Association("Tags").Append(&tag)
}
