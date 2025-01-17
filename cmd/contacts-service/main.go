package main

import (
	"log"
	"net/http"

	"github.com/Horqu/zkp-communicator-backend/internal/contacts"
	"github.com/Horqu/zkp-communicator-backend/internal/db"
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
			"message": "contacts-service pong",
		})
	})

	// AddContactHandler creates a new contact
	r.POST("/contacts", contacts.AddContactHandler(conn))

	// RemoveContactHandler deletes an existing contact
	r.DELETE("/contacts/:contactID", contacts.RemoveContactHandler(conn))

	// VerifyContactStatusHandler gets the status of an existing contact
	r.GET("/contacts/:contactID/status", contacts.VerifyContactStatusHandler(conn))

	r.Run(":8084")
}
