package queries

import (
	"errors"

	"github.com/Horqu/zkp-communicator-backend/cmd/database/models"
	"gorm.io/gorm"
)

func AddMessage(db *gorm.DB, usernameSender, usernameRecipient, C1ForSender, contentForSender, C1ForRecipient, contentForRecipient string) error {
	var sender models.User
	if err := db.Where("username = ?", usernameSender).First(&sender).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("sender not found")
		}
		return err
	}

	var recipient models.User
	if err := db.Where("username = ?", usernameRecipient).First(&recipient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("recipient not found")
		}
		return err
	}

	messageForSender := models.Message{
		ReaderID:    sender.ID,
		SenderID:    sender.ID,
		RecipientID: recipient.ID,
		C1:          C1ForSender,
		Content:     contentForSender,
	}

	messageForRecipient := models.Message{
		ReaderID:    recipient.ID,
		SenderID:    sender.ID,
		RecipientID: recipient.ID,
		C1:          C1ForRecipient,
		Content:     contentForRecipient,
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&messageForSender).Error; err != nil {
			return err
		}
		if err := tx.Create(&messageForRecipient).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func GetMessagesBetweenUsers(db *gorm.DB, readerUsername, otherUsername string) ([]models.Message, error) {
	var reader models.User
	if err := db.Where("username = ?", readerUsername).First(&reader).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reader not found")
		}
		return nil, err
	}

	var otherUser models.User
	if err := db.Where("username = ?", otherUsername).First(&otherUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("other user not found")
		}
		return nil, err
	}

	var messages []models.Message
	if err := db.Where("reader_id = ? AND (sender_id = ? OR recipient_id = ?)", reader.ID, otherUser.ID, otherUser.ID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
