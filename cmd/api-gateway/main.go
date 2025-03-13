package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	sessions = map[string]time.Time{} // Change to username -> pair of token and expiration time
)

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

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
			// Generujemy public key i challenge
			publicKey = "public_key_for_" + message.Data
			challenge := "solve_this_challenge"
			resp := Response{
				Command: "challenge",
				Data:    fmt.Sprintf("%s|%s", publicKey, challenge),
			}
			sendJSON(conn, resp)

		case "solve":
			// Weryfikacja rozwiązania
			if message.Data == "correct_solution" {
				authToken := "auth_token_xxx"
				// Ustawiamy wygaśnięcie np. za 2 minuty
				sessions[authToken] = time.Now().Add(2 * time.Minute)

				resp := Response{Command: "auth", Data: authToken}
				sendJSON(conn, resp)
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid solution"))
			}

		case "ping":
			// Sprawdzamy, czy token jest ważny
			exp, ok := sessions[message.Data]
			if !ok || time.Now().After(exp) {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid token"))
				continue
			}
			// Odświeżamy ważność tokenu
			sessions[message.Data] = time.Now().Add(2 * time.Minute)

			resp := Response{Command: "pong", Data: "token_extended"}
			sendJSON(conn, resp)

		default:
			conn.WriteMessage(websocket.TextMessage, []byte("Unknown command"))
		}
	}
}

func sendJSON(conn *websocket.Conn, resp Response) {
	data, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, data)
}

func main() {
	router := gin.Default()

	router.GET("/ws", wsHandler)

	router.Run(":8080")
}
