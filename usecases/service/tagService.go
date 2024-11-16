package service

import (
	"miw/entities"
	"miw/usecases/repository"
)

type TagUseCase interface {
	CreateTag(tag *entities.Tag) error
	GetTagById(id uint) (*entities.Tag, error)
}

type TagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(tag *entities.Tag) error {
	return s.repo.CreateTag(tag)
}

func (s *TagService) GetTagById(id uint) (*entities.Tag, error) {
	return s.repo.GetTagById(id)
}



