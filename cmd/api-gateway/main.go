package main

import (
	"encoding/json"
	"net/http"

	// "github.com/Horqu/zkp-communicator-backend/internal/auth"

	// "github.com/Horqu/zkp-communicator-backend/internal/messaging"
	// "github.com/Horqu/zkp-communicator-backend/internal/zkp"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type Response struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var authToken string
	var publicKey string

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		switch message.Command {
		case "login":
			// Generate public key and challenge
			publicKey = "public_key_for_" + message.Data
			challenge := "solve_this_challenge"
			response := Response{Command: "challenge", Data: publicKey + "|" + challenge}
			responseMsg, _ := json.Marshal(response)
			conn.WriteMessage(websocket.TextMessage, responseMsg)
		case "solve":
			// Verify solution
			if message.Data == "correct_solution" {
				authToken = "auth_token"
				response := Response{Command: "auth", Data: authToken}
				responseMsg, _ := json.Marshal(response)
				conn.WriteMessage(websocket.TextMessage, responseMsg)
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid solution"))
			}
		case "ping":
			// Extend token validity
			if message.Data == authToken {
				response := Response{Command: "pong", Data: "token_extended"}
				responseMsg, _ := json.Marshal(response)
				conn.WriteMessage(websocket.TextMessage, responseMsg)
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid token"))
			}
		default:
			conn.WriteMessage(websocket.TextMessage, []byte("Unknown command"))
		}
	}
}

func main() {
	router := gin.Default()

	router.GET("/ws", wsHandler)

	router.Run(":8080")
}
