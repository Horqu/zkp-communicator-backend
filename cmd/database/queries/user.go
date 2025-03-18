package queries

import (
	"errors"

	"github.com/Horqu/zkp-communicator-backend/cmd/database/models"
	"gorm.io/gorm"
)

// AddUser dodaje nowego użytkownika do bazy danych
func AddUser(db *gorm.DB, username, publicKey string) error {
	newUser := models.User{
		Username:  username,
		PublicKey: publicKey,
	}
	if err := db.Create(&newUser).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("user already exists")
		}
		return err
	}

	return nil
}

// GetPublicKeyByUsername zwraca klucz publiczny użytkownika na podstawie nazwy użytkownika
func GetPublicKeyByUsername(db *gorm.DB, username string) (string, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	return user.PublicKey, nil
}
