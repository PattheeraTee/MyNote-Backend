package service

import (
	"miw/entities"
	"miw/usecases/repository"
	"fmt"
)

type TagUseCase interface {
	CreateTag(tag *entities.Tag) error
	GetTagById(tagID, userID uint) (*entities.Tag, error)
	UpdateTagName(tagID, userID uint, newName string) error
	DeleteTag(tagID, userID uint) error
}

type TagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// CreateTag: สร้าง Tag พร้อมตรวจสอบว่า User เป็นเจ้าของ
func (s *TagService) CreateTag(tag *entities.Tag) error {
	return s.repo.CreateTag(tag)
}

// GetTagById: ดึง Tag ตาม ID และตรวจสอบ UserID
func (s *TagService) GetTagById(tagID, userID uint) (*entities.Tag, error) {
	tag, err := s.repo.GetTagById(tagID)
	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า Tag เป็นของ User นี้หรือไม่
	if tag.UserID != userID {
		return nil, fmt.Errorf("tag not found or does not belong to this user")
	}

	return tag, nil
}

// UpdateTagName: แก้ไขชื่อ Tag โดยต้องเป็นเจ้าของเท่านั้น
func (s *TagService) UpdateTagName(tagID, userID uint, newName string) error {
	// ตรวจสอบว่าผู้ใช้เป็นเจ้าของแท็กก่อนอัปเดต
	tag, err := s.GetTagById(tagID, userID)
	if err != nil {
		return err
	}

	return s.repo.UpdateTagName(tag.TagID, userID, newName)
}

// DeleteTag: ลบ Tag โดยต้องเป็นเจ้าของเท่านั้น
func (s *TagService) DeleteTag(tagID, userID uint) error {
	// ตรวจสอบว่าผู้ใช้เป็นเจ้าของแท็กก่อนลบ
	tag, err := s.GetTagById(tagID, userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteTag(tag.TagID, userID)
}
