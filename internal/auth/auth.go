package auth

import (
	"net/http"

	"github.com/Horqu/zkp-communicator-backend/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHandler(conn *gorm.DB) gin.HandlerFunc {
	// TODO: Fix registration of a new user
	return func(c *gin.Context) {
		var req struct {
			Username  string `json:"username"`
			PublicKey string `json:"public_key"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		user := db.User{
			Username:  req.Username,
			PublicKey: req.PublicKey,
		}
		if err := conn.Create(&user).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "registered", "user_id": user.ID})
	}
}

func LoginHandler(c *gin.Context) {
	// TODO: Implement user login
	c.JSON(http.StatusOK, gin.H{"token": "generated_jwt_token"})
}

func LogoutHandler(c *gin.Context) {
	// TODO: Implement user logout
	c.JSON(http.StatusOK, gin.H{"status": "logged_out"})
}
