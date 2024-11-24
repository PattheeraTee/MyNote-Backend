package gormRepository

import (
	"miw/entities"
	"fmt"
	"gorm.io/gorm"
	"time"
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
    // if err := r.db.Where("user_id = ?", userID).
	if err := r.db.Where("user_id = ? AND deleted_at = ?", userID, "").
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
    // ใช้เวลาปัจจุบันในรูปแบบ string
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    
    // อัปเดตฟิลด์ DeletedAt ด้วยเวลาปัจจุบัน
    result := r.db.Model(&entities.Note{}).Where("note_id = ? AND deleted_at = ?", id, "").Update("deleted_at", currentTime)

    // ตรวจสอบว่าพบโน้ตหรือไม่
    if result.RowsAffected == 0 {
        return fmt.Errorf("note with ID %d not found or already deleted", id)
    }

    if result.Error != nil {
        return fmt.Errorf("failed to soft delete note with ID %d: %v", id, result.Error)
    }

    return nil
}



func (r *GormNoteRepository) RestoreNoteById(id uint) error {
    // ใช้คำสั่ง Unscoped() เพื่ออัปเดต DeletedAt ให้เป็น nil
    if err := r.db.Unscoped().Model(&entities.Note{}).Where("note_id = ?", id).Update("deleted_at", "").Error; err != nil {
        return fmt.Errorf("failed to restore note with ID %d: %v", id, err)
    }
    return nil
}


func (r *GormNoteRepository) AddTagToNote(noteID uint, tagID uint) error {
    // ดึงข้อมูลโน้ตที่ต้องการเพิ่มแท็ก
    var note entities.Note
    if err := r.db.Preload("Tags").First(&note, noteID).Error; err != nil {
        return fmt.Errorf("note not found: %v", err)
    }

    // ดึงข้อมูลแท็กที่ต้องการเพิ่ม
    var tag entities.Tag
    if err := r.db.First(&tag, tagID).Error; err != nil {
        return fmt.Errorf("tag not found: %v", err)
    }

    // ตรวจสอบว่าแท็กนี้มีอยู่ในโน้ตแล้วหรือยัง
    for _, existingTag := range note.Tags {
        if existingTag.TagID == tagID {
            return fmt.Errorf("tag with ID %d already exists in the note", tagID)
        }
    }

    // เพิ่มแท็กเข้าไปในโน้ต
    if err := r.db.Model(&note).Association("Tags").Append(&tag); err != nil {
        return fmt.Errorf("failed to add tag to note: %v", err)
    }

    return nil
}


func (r *GormNoteRepository) RemoveTagFromNote(noteID uint, tagID uint) error {
    // Fetch the note with the given ID
    var note entities.Note
    if err := r.db.Preload("Tags").First(&note, noteID).Error; err != nil {
        return err
    }

    // Fetch the tag with the given ID
    var tag entities.Tag
    if err := r.db.First(&tag, tagID).Error; err != nil {
        return err
    }

    // Remove the tag from the note's association
    if err := r.db.Model(&note).Association("Tags").Delete(&tag); err != nil {
        return err
    }

    return nil
}
