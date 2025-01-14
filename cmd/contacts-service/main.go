package main

import (
	"net/http"

	"github.com/Horqu/zkp-communicator-backend/internal/contacts"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Ping endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "contacts-service pong",
		})
	})

	// Dodawanie kontaktu
	r.POST("/contacts", contacts.AddContactHandler)

	// Usuwanie kontaktu
	r.DELETE("/contacts/:contactID", contacts.RemoveContactHandler)

	// Weryfikacja statusu kontaktu
	r.GET("/contacts/:contactID/status", contacts.VerifyContactStatusHandler)

	r.Run(":8084")
}
