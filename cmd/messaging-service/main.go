package main

import (
	"log"
	"net/http"

	"github.com/Horqu/zkp-communicator-backend/internal/db"
	"github.com/Horqu/zkp-communicator-backend/internal/messaging"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "messaging-service pong",
		})
	})

	// Obsługa przesyłania zaszyfrowanych wiadomości
	r.POST("/messages", messaging.SendMessageHandler(conn))

	// Obsługa odbierania zaszyfrowanych wiadomości
	r.GET("/messages", messaging.ReceiveMessagesHandler(conn))

	// Zarządzanie historią rozmów
	r.GET("/messages/history", messaging.GetHistoryHandler)

	// Tworzenie grupy
	r.POST("/groups", messaging.CreateGroupHandler)

	// Generowanie kluczy grupowych (przykładowa operacja)
	r.POST("/groups/:groupID/keys", messaging.GenerateGroupKeysHandler)

	// Weryfikacja kluczy grupowych (przykładowa operacja)
	r.POST("/groups/:groupID/keys/verify", messaging.VerifyGroupKeysHandler)

	// Usuwanie grupy i jej zawartości
	r.DELETE("/groups/:groupID", messaging.DeleteGroupHandler)

	r.Run(":8083")
}
