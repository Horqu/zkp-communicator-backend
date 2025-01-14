package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Horqu/zkp-communicator-backend/internal/zkp"
)

func main() {
	router := gin.Default()

	// Endpoint for ZKP proof generation
	router.POST("/generate-proof", zkp.GenerateProofHandler)

	// Endpoint for ZKP proof verification
	router.POST("/verify-proof", zkp.VerifyProofHandler)

	// Ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "zkp-service pong",
		})
	})

	router.Run(":8081") // Run on port 8081
}
