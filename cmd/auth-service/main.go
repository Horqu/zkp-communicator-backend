package main

import (
	"net/http"

	"github.com/Horqu/zkp-communicator-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "auth-service pong",
		})
	})

	router.POST("/register", auth.RegisterHandler)

	router.POST("/login", auth.LoginHandler)

	router.POST("/logout", auth.LogoutHandler)

	router.Run(":8082")
}
