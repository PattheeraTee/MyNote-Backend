package handler

import (
	"miw/entities"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	result := db.Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func Login(db *gorm.DB, user *entities.User) (string, error) {
	selectedUser := new(entities.User)
	result := db.Where("email = ?", user.Email).First(selectedUser)

	if result.Error != nil {
		return "", result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(selectedUser.Password), []byte(user.Password))

	if err != nil {
		return "", err
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = selectedUser.UserID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return t, nil
}

func GetUser(db *gorm.DB, id uint) (*entities.User) {
	var user entities.User
	result := db.First(&user, id)

	if result.Error != nil {
		return nil
	}

	return &user
}

func UpdateUsername(db *gorm.DB, user *entities.User) error {
	
	result := db.Model(&user).Updates(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdatePassword(db *gorm.DB, user *entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	result := db.Model(&user).Updates(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}