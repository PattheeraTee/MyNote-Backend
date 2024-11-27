package repository

import (
	"miw/entities"
)

type TagRepository interface {
	CreateTag(tag *entities.Tag) error
	GetTagById(id uint) (*entities.Tag, error) 
	GetTagsByUser(userID uint) ([]entities.Tag, error)
	UpdateTagName(tagID, userID uint, newName string) error
	DeleteTag(tagID, userID uint) error
}
