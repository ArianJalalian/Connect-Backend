package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(user dto.User) (models.User, error) {
	if err := database.DB.Create(&user).Error; err != nil {
		fmt.Println("error is: ", err)
		return models.User{}, err
	}
	uModel, err := GetUserByEmail(user.Email)
	return *uModel, err
}

func GetUser(email, password string) (models.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	fmt.Println("user: ", user)
	if user == nil {
		return models.User{}, nil
	}
	fmt.Println(user.Password, password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserById(id uint64) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
