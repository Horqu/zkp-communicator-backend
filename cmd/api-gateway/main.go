package main

import (
	"net/http"

	// "github.com/Horqu/zkp-communicator-backend/internal/auth"

	// "github.com/Horqu/zkp-communicator-backend/internal/messaging"
	// "github.com/Horqu/zkp-communicator-backend/internal/zkp"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// // Rejestracja użytkownika ZKP
	// router.POST("/register", auth.RegisterHandler)

	// // Logowanie ZKP
	// router.POST("/login", auth.LoginHandler)

	// // Wylogowanie użytkownika
	// router.POST("/logout", auth.LogoutHandler)

	// // Dodawanie kontaktu
	// router.POST("/contacts", auth.AddContactHandler)

	// // Usuwanie kontaktu
	// router.DELETE("/contacts/:id", auth.RemoveContactHandler)

	// // Tworzenie rozmowy grupowej
	// router.POST("/groups", messaging.CreateGroupHandler)

	// // Usuwanie rozmowy grupowej
	// router.DELETE("/groups/:id", messaging.DeleteGroupHandler)

	// // Wysyłanie wiadomości
	// router.POST("/messages", messaging.SendMessageHandler)

	// // Odbieranie wiadomości
	// router.GET("/messages", messaging.ReceiveMessagesHandler)

	// // Weryfikacja wieku
	// router.POST("/verify-age", zkp.VerifyAgeHandler)

	// // Szyfrowana komunikacja middleware
	// router.Use(encryption.EncryptMiddleware(), encryption.DecryptMiddleware())

	// Monitorowanie systemu (Prometheus)
	// Endpointy Prometheus są zazwyczaj obsługiwane osobno

	router.Run(":8080")
}
