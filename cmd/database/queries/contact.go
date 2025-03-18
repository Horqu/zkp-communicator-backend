package queries

import (
	"errors"

	"github.com/Horqu/zkp-communicator-backend/cmd/database/models"
	"gorm.io/gorm"
)

func AddContact(db *gorm.DB, username1, username2 string) error {
	var user1, user2 models.User
	if err := db.Where("username = ?", username1).First(&user1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user1 not found")
		}
		return err
	}
	if err := db.Where("username = ?", username2).First(&user2).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user2 not found")
		}
		return err
	}

	var contact1, contact2 models.Contact
	err1 := db.Where("user_id = ? AND contact_id = ?", user1.ID, user2.ID).First(&contact1).Error
	err2 := db.Where("user_id = ? AND contact_id = ?", user2.ID, user1.ID).First(&contact2).Error

	if errors.Is(err1, gorm.ErrRecordNotFound) && errors.Is(err2, gorm.ErrRecordNotFound) {
		contact1 = models.Contact{
			UserID:    user1.ID,
			ContactID: user2.ID,
			Status:    "request_sent",
		}
		contact2 = models.Contact{
			UserID:    user2.ID,
			ContactID: user1.ID,
			Status:    "pending",
		}
		if err := db.Create(&contact1).Error; err != nil {
			return err
		}
		if err := db.Create(&contact2).Error; err != nil {
			return err
		}
		return nil
	}

	if contact1.Status == "pending" && contact2.Status == "request_sent" {
		if err := db.Model(&contact1).Update("status", "accepted").Error; err != nil {
			return err
		}
		if err := db.Model(&contact2).Update("status", "accepted").Error; err != nil {
			return err
		}
		return nil
	}

	return errors.New("contact already exists or invalid status")
}

func GetContactsByUsername(db *gorm.DB, username string) ([]models.Contact, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var contacts []models.Contact
	if err := db.Where("user_id = ?", user.ID).Preload("Contact").Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}
