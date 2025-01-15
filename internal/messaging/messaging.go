package messaging

import (
	"net/http"
	"strconv"

	"github.com/Horqu/zkp-communicator-backend/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SendMessageHandler obsługuje wysyłanie zaszyfrowanych wiadomości
func SendMessageHandler(conn *gorm.DB) gin.HandlerFunc {
	// TODO: Implementacja szyfrowania, zapisu i wysyłki wiadomości
	return func(c *gin.Context) {
		var req struct {
			SenderID    uint   `json:"sender_id"`
			RecipientID uint   `json:"recipient_id"`
			Content     string `json:"content"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		msg := db.Message{
			SenderID:    req.SenderID,
			RecipientID: req.RecipientID,
			Content:     req.Content,
		}
		if err := conn.Create(&msg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "message_sent"})
	}
}

// ReceiveMessagesHandler obsługuje odbieranie zaszyfrowanych wiadomości
func ReceiveMessagesHandler(conn *gorm.DB) gin.HandlerFunc {
	// TODO: Implementacja pobierania i deszyfrowania wiadomości
	return func(c *gin.Context) {
		user1 := c.Query("user1")
		user2 := c.Query("user2")
		if user1 == "" || user2 == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user1/user2 query params"})
			return
		}

		var messages []db.Message
		// Konwersja ID na uint
		u1, err1 := strconv.ParseUint(user1, 10, 64)
		u2, err2 := strconv.ParseUint(user2, 10, 64)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user IDs"})
			return
		}

		// Znajdź wszystkie wiadomości w obie strony
		if err := conn.Where(
			`(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)`,
			u1, u2, u2, u1,
		).Find(&messages).Error; err != nil && err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}

// GetHistoryHandler obsługuje zarządzanie i pobieranie historii rozmów
func GetHistoryHandler(c *gin.Context) {
	// TODO: Implementacja pobierania historii czatu
	c.JSON(http.StatusOK, gin.H{"history": "chat_history"})
}
