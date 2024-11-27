package gormRepository

import (
	"miw/entities"
	"gorm.io/gorm"
	"fmt"
)

type GormTagRepository struct {
	db *gorm.DB
}

func NewGormTagRepository(db *gorm.DB) *GormTagRepository {
	return &GormTagRepository{db: db}
}

func (r *GormTagRepository) CreateTag(tag *entities.Tag) error {
    // ตรวจสอบว่าชื่อแท็กซ้ำสำหรับ User เดียวกันหรือไม่
    var existingTag entities.Tag
    if err := r.db.Where("user_id = ? AND tag_name = ?", tag.UserID, tag.TagName).First(&existingTag).Error; err == nil {
        return fmt.Errorf("tag name '%s' already exists for this user", tag.TagName)
    }

    // สร้างแท็กใหม่
    if err := r.db.Create(tag).Error; err != nil {
        return fmt.Errorf("failed to create tag: %v", err)
    }

    return nil
}

func (r *GormTagRepository) GetTagsByUser(userID uint) ([]entities.Tag, error) {
    var tags []entities.Tag
    if err := r.db.Where("user_id = ?", userID).Find(&tags).Error; err != nil {
        return nil, fmt.Errorf("failed to fetch tags: %v", err)
    }
    return tags, nil
}


func (r *GormTagRepository) UpdateTagName(tagID, userID uint, newName string) error {
    // ตรวจสอบว่าแท็กมีอยู่และเป็นของ User นี้หรือไม่
    var tag entities.Tag
    if err := r.db.Where("tag_id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("tag not found or does not belong to this user")
        }
        return fmt.Errorf("error finding tag: %v", err)
    }

    // ตรวจสอบว่าชื่อใหม่ซ้ำหรือไม่
    var existingTag entities.Tag
    if err := r.db.Where("user_id = ? AND tag_name = ?", userID, newName).First(&existingTag).Error; err == nil {
        return fmt.Errorf("tag name '%s' already exists for this user", newName)
    }

    // อัปเดตชื่อแท็ก
    if err := r.db.Model(&tag).Update("tag_name", newName).Error; err != nil {
        return fmt.Errorf("failed to update tag name: %v", err)
    }

    return nil
}



func (r *GormTagRepository) DeleteTag(tagID, userID uint) error {
    // ตรวจสอบว่าแท็กมีอยู่และเป็นของ User นี้หรือไม่
    var tag entities.Tag
    if err := r.db.Where("tag_id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("tag not found or does not belong to this user")
        }
        return fmt.Errorf("error finding tag: %v", err)
    }

    // ลบแท็ก
    if err := r.db.Delete(&tag).Error; err != nil {
        return fmt.Errorf("failed to delete tag: %v", err)
    }

    return nil
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

