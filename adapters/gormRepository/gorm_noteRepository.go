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
	// ใช้ transaction เพื่อความปลอดภัย
	return r.db.Transaction(func(tx *gorm.DB) error {
		// สร้าง Note
		if err := tx.Create(note).Error; err != nil {
			return fmt.Errorf("failed to create note: %v", err)
		}

		// เพิ่ม ToDo Items ถ้ามี
		if len(note.TodoItems) > 0 {
            for i := range note.TodoItems {
                note.TodoItems[i].NoteID = note.NoteID // ตั้ง NoteID ให้กับ ToDo Item
                note.TodoItems[i].ID = 0              // รีเซ็ตค่า ID ให้เป็น 0 เพื่อให้ฐานข้อมูลจัดการ Auto Increment
            }
            if err := tx.Create(&note.TodoItems).Error; err != nil {
                return fmt.Errorf("failed to create todo items: %v", err)
            }
        }        
		return nil
	})
}


func (r *GormNoteRepository) UpdateNote(note *entities.Note) error {
	// อัปเดต Note
	if err := r.db.Save(note).Error; err != nil {
		return fmt.Errorf("failed to update note: %v", err)
	}

	// ลบ TodoItems เก่าทั้งหมดก่อนเพิ่มใหม่ (หรือใช้ Merge หากต้องการอัปเดตเฉพาะ)
	if err := r.db.Where("note_id = ?", note.NoteID).Delete(&entities.ToDo{}).Error; err != nil {
		return fmt.Errorf("failed to delete old todo items: %v", err)
	}

	// เพิ่ม TodoItems ใหม่
	for i := range note.TodoItems {
		note.TodoItems[i].NoteID = note.NoteID
		if err := r.db.Create(&note.TodoItems[i]).Error; err != nil {
			return fmt.Errorf("failed to create new todo item: %v", err)
		}
	}
	return nil
}


func (r *GormNoteRepository) GetAllNoteByUserId(userID uint) ([]entities.Note, error) {
	var notes []entities.Note
	if err := r.db.Where("user_id = ? AND deleted_at = ?", userID, "").
        Preload("Tags", func(db *gorm.DB) *gorm.DB {
            return db.Select("tag_id, tag_name") // ไม่ดึง Notes ใน Tags
        }).
        Preload("Reminder").
        Preload("Event").
        Preload("TodoItems"). // เพิ่มการโหลด TodoItems
        Find(&notes).Error; err != nil {
        return nil, err
    }
	return notes, nil
}

func (r *GormNoteRepository) GetNoteById(id uint) (*entities.Note, error) {
	var note entities.Note
	if err := r.db.Where("note_id = ?", id).
	Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("tag_id, tag_name") // ไม่ดึง Notes ใน Tags
	}).
	Preload("Reminder").
	Preload("Event").
    Preload("TodoItems"). // เพิ่มการโหลด TodoItems
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
    // ตรวจสอบว่ามี Note ที่ตรงกับ ID หรือไม่
    var note entities.Note
    if err := r.db.Unscoped().Where("note_id = ?", id).First(&note).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("note with ID %d not found", id)
        }
        return fmt.Errorf("failed to check note with ID %d: %v", id, err)
    }

    // ใช้คำสั่ง Unscoped() เพื่ออัปเดต DeletedAt ให้เป็น nil
    if err := r.db.Unscoped().Model(&entities.Note{}).Where("note_id = ?", id).Update("deleted_at", "").Error; err != nil {
        return fmt.Errorf("failed to restore note with ID %d: %v", id, err)
    }

    return nil
}



func (r *GormNoteRepository) AddTagToNote(noteID uint, tagID uint, userID uint) error {
    // ตรวจสอบว่า Note เป็นของ User หรือไม่
    var note entities.Note
    if err := r.db.Where("note_id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
        return fmt.Errorf("note not found or does not belong to the user")
    }

    // ตรวจสอบว่า Tag เป็นของ User หรือไม่
    var tag entities.Tag
    if err := r.db.Where("tag_id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
        return fmt.Errorf("tag not found or does not belong to the user")
    }

    // เพิ่ม Tag เข้า Note
    if err := r.db.Model(&note).Association("Tags").Append(&tag); err != nil {
        return fmt.Errorf("failed to add tag to note: %v", err)
    }

    return nil
}


func (r *GormNoteRepository) RemoveTagFromNote(noteID uint, tagID uint, userID uint) error {
    // ตรวจสอบว่า Note เป็นของ User หรือไม่
    var note entities.Note
    if err := r.db.Where("note_id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
        return fmt.Errorf("note not found or does not belong to the user")
    }

    // ตรวจสอบว่า Tag เป็นของ User หรือไม่
    var tag entities.Tag
    if err := r.db.Where("tag_id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
        return fmt.Errorf("tag not found or does not belong to the user")
    }

    // ลบ Tag ออกจาก Note
    if err := r.db.Model(&note).Association("Tags").Delete(&tag); err != nil {
        return fmt.Errorf("failed to remove tag from note: %v", err)
    }

    return nil
}

func (r *GormNoteRepository) GetNoteByIdAndUser(noteID uint, userID uint) (*entities.Note, error) {
	var note entities.Note
	if err := r.db.Where("note_id = ? AND user_id = ?", noteID, userID).
		Preload("Tags").
		Preload("Reminder").
		Preload("Event").
        Preload("TodoItems").
		First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("note not found or does not belong to the user")
		}
		return nil, err
	}
	return &note, nil
}


