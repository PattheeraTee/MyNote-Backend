package repository

import (
	"miw/entities"
)

type TagRepository interface {
	CreateTag(tag *entities.Tag) error
	GetTagById(id uint) (*entities.Tag, error)
}
