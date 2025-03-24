package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "github.com/Horqu/zkp-communicator-backend/cmd/database"
	"github.com/Horqu/zkp-communicator-backend/cmd/database/queries"
	"github.com/Horqu/zkp-communicator-backend/cmd/internal"
	wsresponses "github.com/Horqu/zkp-communicator-backend/cmd/ws-responses"

	// "github.com/Horqu/zkp-communicator-backend/cmd/encryption"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	sessions   = map[string]time.Time{} // Change to username -> pair of token and expiration time
	challenges = map[string]string{}    // Change to username -> challenge
)

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var message internal.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		switch message.Command {
		case internal.MessageLoginButtom:
			resp := wsresponses.ResponseLoginPage()
			sendJSON(conn, resp)
		case internal.MessageRegister:
			var dataMap map[string]string
			err := json.Unmarshal([]byte(message.Data), &dataMap)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format"))
				continue
			}

			username := dataMap["username"]
			publicKey := dataMap["publicKey"]
			log.Printf("Registering user: username=%s, publicKey=%s\n", username, publicKey)

			err = queries.AddUser(db.GetDBInstance(), username, publicKey)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to register user"))
				return
			}

			resp := wsresponses.ResponseRegisterSuccess()
			sendJSON(conn, resp)

		case internal.MessageLogin:
			var dataMap map[string]string
			err := json.Unmarshal([]byte(message.Data), &dataMap)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format"))
				continue
			}

			username := dataMap["username"]
			method := dataMap["method"]
			log.Printf("Logging in user: username=%s, method=%s\n", username, method)

			publicKey, err := queries.GetPublicKeyByUsername(db.GetDBInstance(), username)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			if method == "methodA" {
				resp := wsresponses.ResponseLoginSuccess(publicKey)
				sendJSON(conn, resp)
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid method"))
				continue
			}
		case internal.MessageSolve:
			var dataMap map[string]string
			err := json.Unmarshal([]byte(message.Data), &dataMap)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format"))
				continue
			}

			username := dataMap["username"]
			log.Printf("Resolving user: username=%s\n", username)

			publicKey, err := queries.GetPublicKeyByUsername(db.GetDBInstance(), username)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			friendList, err := queries.GetContactsByUsername(db.GetDBInstance(), username)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to get contacts"))
				continue
			}

			resp := wsresponses.ResponseSolveSuccess(publicKey, friendList)
			sendJSON(conn, resp)

		case internal.MessageAddFriend:
			var dataMap map[string]string
			err := json.Unmarshal([]byte(message.Data), &dataMap)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format"))
				continue
			}

			username := dataMap["username"]
			friend := dataMap["friend"]
			log.Printf("Adding friend: username=%s, friend=%s\n", username, friend)

			publicKey, err := queries.GetPublicKeyByUsername(db.GetDBInstance(), friend)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Friend not found"))
				continue
			}

			err = queries.AddContact(db.GetDBInstance(), username, friend)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to add friend"))
				continue
			}

			friendList, err := queries.GetContactsByUsername(db.GetDBInstance(), username)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to get contacts"))
				continue
			}

			resp := wsresponses.ResponseSolveSuccess(publicKey, friendList)
			sendJSON(conn, resp)
		case internal.MessageSelectChat:
			var dataMap map[string]string
			err := json.Unmarshal([]byte(message.Data), &dataMap)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format"))
				continue
			}

			username := dataMap["username"]
			friend := dataMap["friend"]
			log.Printf("Selecting chat: username=%s, friend=%s\n", username, friend)

			type SimplifiedMessage struct {
				SenderUsername    string    `json:"senderUsername"`
				RecipientUsername string    `json:"recipientUsername"`
				C1                string    `json:"c1"`
				Content           string    `json:"content"`
				CreatedAt         time.Time `json:"createdAt"`
			}

			friendPublicKey, err := queries.GetPublicKeyByUsername(db.GetDBInstance(), friend)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Friend not found"))
				continue
			}

			chatForUser, err := queries.GetMessagesBetweenUsers(db.GetDBInstance(), username, friend)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to get chat"))
				continue
			}

			// Przekształć dane w listę SimplifiedMessage
			var simplifiedMessages []SimplifiedMessage
			for _, message := range chatForUser {
				senderUsername, _ := queries.GetUsernameByUserID(db.GetDBInstance(), message.SenderID)
				recipientUsername, _ := queries.GetUsernameByUserID(db.GetDBInstance(), message.RecipientID)

				simplifiedMessages = append(simplifiedMessages, SimplifiedMessage{
					SenderUsername:    senderUsername,
					RecipientUsername: recipientUsername,
					C1:                message.C1,
					Content:           message.Content,
					CreatedAt:         message.CreatedAt,
				})
			}

			// Tworzymy strukturę odpowiedzi z friendPublicKey
			responseData := struct {
				FriendPublicKey string              `json:"friendPublicKey"`
				Messages        []SimplifiedMessage `json:"messages"`
			}{
				FriendPublicKey: friendPublicKey,
				Messages:        simplifiedMessages,
			}

			// Serializuj dane do JSON
			data, err := json.Marshal(responseData)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize chat data"))
				continue
			}

			// Wyślij dane jako odpowiedź
			resp := internal.Response{
				Command: internal.ResponseSelectChat,
				Data:    string(data),
			}
			sendJSON(conn, resp)

		case internal.MessageSendMessage:
			var dataMap map[string]string
			log.Printf("Received message: %s\n", message.Data)
			err := json.Unmarshal([]byte(message.Data), &dataMap)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format"))
				continue
			}

			username := dataMap["username"]
			friend := dataMap["friend"]
			c1user := dataMap["c1user"]
			contentuser := dataMap["contentuser"]
			c1friend := dataMap["c1friend"]
			contentfriend := dataMap["contentfriend"]
			log.Printf("Sending message: username=%s, friend=%s, c1user=%s, contentuser=%s, c1friend=%s, contentfriend=%s\n", username, friend, c1user, contentuser, c1friend, contentfriend)

			err = queries.AddMessage(db.GetDBInstance(), username, friend, c1user, contentuser, c1friend, contentfriend)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to send message"))
				continue
			}

			type SimplifiedMessage struct {
				SenderUsername    string    `json:"senderUsername"`
				RecipientUsername string    `json:"recipientUsername"`
				C1                string    `json:"c1"`
				Content           string    `json:"content"`
				CreatedAt         time.Time `json:"createdAt"`
			}

			friendPublicKey, err := queries.GetPublicKeyByUsername(db.GetDBInstance(), friend)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Friend not found"))
				continue
			}

			chatForUser, err := queries.GetMessagesBetweenUsers(db.GetDBInstance(), username, friend)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to get chat"))
				continue
			}

			queries.GetUsernameByUserID(db.GetDBInstance(), 1)
			// Przekształć dane w listę SimplifiedMessage
			var simplifiedMessages []SimplifiedMessage
			for _, message := range chatForUser {
				senderUsername, _ := queries.GetUsernameByUserID(db.GetDBInstance(), message.SenderID)
				recipientUsername, _ := queries.GetUsernameByUserID(db.GetDBInstance(), message.RecipientID)

				simplifiedMessages = append(simplifiedMessages, SimplifiedMessage{
					SenderUsername:    senderUsername,
					RecipientUsername: recipientUsername,
					C1:                message.C1,
					Content:           message.Content,
					CreatedAt:         message.CreatedAt,
				})
			}

			// Tworzymy strukturę odpowiedzi z friendPublicKey
			responseData := struct {
				FriendPublicKey string              `json:"friendPublicKey"`
				Messages        []SimplifiedMessage `json:"messages"`
			}{
				FriendPublicKey: friendPublicKey,
				Messages:        simplifiedMessages,
			}

			// Serializuj dane do JSON
			data, err := json.Marshal(responseData)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize chat data"))
				continue
			}

			// Wyślij dane jako odpowiedź
			resp := internal.Response{
				Command: internal.ResponseSelectChat,
				Data:    string(data),
			}
			log.Printf("Sending data: %s", data)

			sendJSON(conn, resp)

		case "ping":
			// Sprawdzamy, czy token jest ważny
			exp, ok := sessions[message.Data]
			if !ok || time.Now().After(exp) {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid token"))
				continue
			}
			// Odświeżamy ważność tokenu
			sessions[message.Data] = time.Now().Add(2 * time.Minute)

			resp := internal.Response{Command: "pong", Data: "token_extended"}
			sendJSON(conn, resp)

		default:
			conn.WriteMessage(websocket.TextMessage, []byte("Unknown command"))
		}
	}
}

func sendJSON(conn *websocket.Conn, resp internal.Response) {
	data, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, data)
}

func main() {
	db.GetDBInstance()

	router := gin.Default()

	router.GET("/ws", wsHandler)

	router.Run(":8080")
}
