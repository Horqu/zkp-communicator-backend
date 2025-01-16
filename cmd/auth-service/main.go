package main

import (
	"log"
	"net/http"

	"github.com/Horqu/zkp-communicator-backend/internal/auth"
	"github.com/Horqu/zkp-communicator-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Connect to the database
	conn, err := db.ConnectGORM()
	if err != nil {
		log.Fatal("DB connection failed: ", err)
	}

	// Auto-migrate the database
	if err := db.AutoMigrateAll(conn); err != nil {
		log.Fatal("DB migration failed: ", err)
	}

	log.Default().Println("DB connection successful")

	// Ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "auth-service pong",
		})
	})

	router.POST("/register", auth.RegisterHandler(conn))

	router.POST("/login", auth.LoginHandler(conn))

	router.POST("/logout", auth.LogoutHandler(conn))

	router.Run(":8082")
}
