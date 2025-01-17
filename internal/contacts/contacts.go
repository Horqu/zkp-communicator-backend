package contacts

import (
	"net/http"

	"github.com/Horqu/zkp-communicator-backend/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddContactHandler creates a new contact
func AddContactHandler(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID        uint   `json:"user_id"`
			ContactUserID uint   `json:"contact_user_id"`
			Status        string `json:"status"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		newContact := db.Contact{
			UserID:        req.UserID,
			ContactUserID: req.ContactUserID,
			Status:        req.Status,
		}
		if err := conn.Create(&newContact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "contact_added", "contact_id": newContact.ID})
	}
}

// RemoveContactHandler deletes an existing contact
func RemoveContactHandler(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		contactID := c.Param("contactID")
		var contact db.Contact
		if err := conn.Where("contact_id = ?", contactID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		if err := conn.Delete(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove contact"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "contact_removed", "contact_id": contactID})
	}
}

// VerifyContactStatusHandler gets the status of an existing contact
func VerifyContactStatusHandler(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		contactID := c.Param("contactID")
		var contact db.Contact
		if err := conn.Where("contact_id = ?", contactID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":          contact.Status,
			"user_id":         contact.UserID,
			"contact_user_id": contact.ContactUserID,
		})
	}
}
