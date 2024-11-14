package gormRepository

import (
	"gorm.io/gorm"
	"miw/entities"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) GetUserById(id uint) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}