package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "github.com/Horqu/zkp-communicator-backend/cmd/database"
	"github.com/Horqu/zkp-communicator-backend/cmd/database/queries"
	"github.com/Horqu/zkp-communicator-backend/cmd/encryption"
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

	activeUsers                     = map[string]internal.ActiveUser{}
	activeSchnorrChallenges         = map[string]internal.SchnorrChallengeToSave{}
	activeFeigeFiatShamirChallenges = map[string]internal.FeigeFiatShamirChallengeToSave{}
	activeSigmaChallenges           = map[string]internal.SigmaChallengeToSave{}
)

func startUserSessionChecker() {
	ticker := time.NewTicker(1 * time.Second) // Timer, który działa co sekundę
	defer ticker.Stop()

	for {
		<-ticker.C
		now := time.Now()

		for username, user := range activeUsers {
			// Sprawdź, czy czas sesji użytkownika minął
			if now.After(user.Expiry) {
				log.Printf("Session expired for user: %s\n", username)
				user.WsConnection.WriteMessage(websocket.TextMessage, []byte("Session expired"))
				user.WsConnection.Close()     // Zamknij połączenie WebSocket
				delete(activeUsers, username) // Usuń użytkownika z mapy
				continue
			}

			if user.SelectedFriend != "" {

				friend := user.SelectedFriend
				log.Printf("Checking chat for user: %s with friend: %s\n", username, friend)
				conn := user.WsConnection

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

				responseData := struct {
					FriendPublicKey string              `json:"friendPublicKey"`
					Messages        []SimplifiedMessage `json:"messages"`
				}{
					FriendPublicKey: friendPublicKey,
					Messages:        simplifiedMessages,
				}

				data, err := json.Marshal(responseData)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize chat data"))
					continue
				}

				resp := internal.Response{
					Command: internal.ResponseSelectChat,
					Data:    string(data),
				}

				sendJSON(conn, resp)
				continue
			}

			err := user.WsConnection.WriteMessage(websocket.TextMessage, []byte("You are still active"))
			if err != nil {
				log.Printf("Failed to send message to user: %s, removing from active users\n", username)
				user.WsConnection.Close()
				delete(activeUsers, username)
			}
		}

		for username, challenge := range activeSchnorrChallenges {
			if time.Now().After(challenge.Expiry) {
				log.Printf("Schnorr challenge expired for user: %s\n", username)
				delete(activeSchnorrChallenges, username)
			}
		}

		for username, challenge := range activeFeigeFiatShamirChallenges {
			if time.Now().After(challenge.Expiry) {
				log.Printf("FeigeFiatShamir challenge expired for user: %s\n", username)
				delete(activeFeigeFiatShamirChallenges, username)
			}
		}

		for username, challenge := range activeSigmaChallenges {
			if time.Now().After(challenge.Expiry) {
				log.Printf("Sigma challenge expired for user: %s\n", username)
				delete(activeSigmaChallenges, username)
			}
		}
	}
}

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

			if method == "Schnorr" {
				e, r, R := encryption.GenerateSchnorrChallenge(publicKey)
				log.Printf("Generated Schnorr challenge: e=%s, r=%s, R=%s\n", e.String(), r.String(), R.String())

				activeSchnorrChallenges[username] = internal.SchnorrChallengeToSave{
					PublicKey: publicKey,
					R:         R,
					E:         e,
					Expiry:    time.Now().Add(2 * time.Minute),
				}

				activeSchnorrToSend := internal.SchnorrChallengeToSend{
					R: r,
					E: e,
				}

				activeSchnorrToSendJson, err := json.Marshal(activeSchnorrToSend)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize challenge"))
					continue
				}

				resp := internal.Response{
					Command: internal.ResponseSchnorrChallenge,
					Data:    string(activeSchnorrToSendJson),
				}

				sendJSON(conn, resp)
			} else if method == "FeigeFiatShamir" {
				N, e := encryption.GenerateFeigeFiatShamirChallenge()
				activeFeigeFiatShamirChallenges[username] = internal.FeigeFiatShamirChallengeToSave{
					N:      N,
					E:      e,
					Expiry: time.Now().Add(2 * time.Minute),
				}

				publicKeyG1Affine, err := encryption.StringToG1Affine(publicKey)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to convert public key"))
					continue
				}

				c1N, encryptedN := encryption.EncryptText(N.String(), &publicKeyG1Affine)

				c1e, encryptedE := encryption.EncryptText(e.String(), &publicKeyG1Affine)

				activeFeigeFiatShamirToSend := internal.FeigeFiatShamirChallengeToSend{
					C1N:        c1N,
					EncryptedN: encryptedN,
					C1e:        c1e,
					EncryptedE: encryptedE,
				}

				activeFeigeFiatShamirToSendJson, err := json.Marshal(activeFeigeFiatShamirToSend)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize challenge"))
					continue
				}

				resp := internal.Response{
					Command: internal.ResponseFFSChallenge,
					Data:    string(activeFeigeFiatShamirToSendJson),
				}

				sendJSON(conn, resp)
			} else if method == "Sigma" {
				e, r, t := encryption.GenerateSigmaChallenge()
				activeSigmaChallenges[username] = internal.SigmaChallengeToSave{
					E:      e,
					T:      t,
					Expiry: time.Now().Add(2 * time.Minute),
				}

				publicKeyG1Affine, err := encryption.StringToG1Affine(publicKey)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to convert public key"))
					continue
				}
				c1e, encryptedE := encryption.EncryptText(e.String(), &publicKeyG1Affine)
				c1r, encryptedR := encryption.EncryptText(r.String(), &publicKeyG1Affine)
				activeSigmaToSend := internal.SigmaChallengeToSend{
					C1e:        c1e,
					EncryptedE: encryptedE,
					C1r:        c1r,
					EncryptedR: encryptedR,
				}
				activeSigmaToSendJson, err := json.Marshal(activeSigmaToSend)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize challenge"))
					continue
				}

				resp := internal.Response{
					Command: internal.ResponseSigmaChallenge,
					Data:    string(activeSigmaToSendJson),
				}

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
			method := dataMap["method"]
			log.Printf("Solving challenge: username=%s, method=%s\n", username, method)

			publicKey, err := queries.GetPublicKeyByUsername(db.GetDBInstance(), username)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			if method == "Schnorr" {
				s := dataMap["s"]
				challenge, exists := activeSchnorrChallenges[username]
				if !exists {
					log.Printf("Challenge not found for user: %s\n", username)
					conn.WriteMessage(websocket.TextMessage, []byte("Challenge not found"))
					continue
				}
				if time.Now().After(challenge.Expiry) {
					log.Printf("Schnorr challenge expired for user: %s\n", username)
					conn.WriteMessage(websocket.TextMessage, []byte("Challenge expired"))
					continue
				}

				sInt, err := encryption.PublicKeyStringToBigInt(s)
				if err != nil {
					log.Printf("Failed to convert s to big.Int: %v\n", err)
					conn.WriteMessage(websocket.TextMessage, []byte("Invalid s value"))
					continue
				}
				R := challenge.R
				e := challenge.E
				log.Printf("Verifying Schnorr proof: R=%s, e=%s, s=%s\n", R.String(), e.String(), sInt.String())
				if !encryption.VerifySchnorrProof(R, e, sInt, publicKey) {
					log.Printf("Schnorr proof verification failed for user: %s\n", username)
					conn.WriteMessage(websocket.TextMessage, []byte("Invalid proof"))
					continue
				}
				delete(activeSchnorrChallenges, username)
				log.Printf("Schnorr proof verified for user: %s\n", username)
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid method"))
				continue
			}

			activeUsers[username] = internal.ActiveUser{
				WsConnection:   conn,
				Expiry:         time.Now().Add(2 * time.Minute),
				PublicKey:      publicKey,
				SelectedFriend: "",
			}
			log.Printf("User %s logged in successfully\n", username)

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

			_, exist := activeUsers[username]
			if !exist {
				conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			if time.Now().After(activeUsers[username].Expiry) {
				conn.WriteMessage(websocket.TextMessage, []byte("Session expired"))
				continue
			}

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

			_, exist := activeUsers[username]
			if !exist {
				conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			if time.Now().After(activeUsers[username].Expiry) {
				conn.WriteMessage(websocket.TextMessage, []byte("Session expired"))
				continue
			}

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

			responseData := struct {
				FriendPublicKey string              `json:"friendPublicKey"`
				Messages        []SimplifiedMessage `json:"messages"`
			}{
				FriendPublicKey: friendPublicKey,
				Messages:        simplifiedMessages,
			}

			data, err := json.Marshal(responseData)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize chat data"))
				continue
			}

			resp := internal.Response{
				Command: internal.ResponseSelectChat,
				Data:    string(data),
			}

			activeUsers[username] = internal.ActiveUser{
				WsConnection:   conn,
				Expiry:         time.Now().Add(2 * time.Minute),
				PublicKey:      activeUsers[username].PublicKey,
				SelectedFriend: friend,
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

			_, exist := activeUsers[username]
			if !exist {
				conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			if time.Now().After(activeUsers[username].Expiry) {
				conn.WriteMessage(websocket.TextMessage, []byte("Session expired"))
				continue
			}

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

			responseData := struct {
				FriendPublicKey string              `json:"friendPublicKey"`
				Messages        []SimplifiedMessage `json:"messages"`
			}{
				FriendPublicKey: friendPublicKey,
				Messages:        simplifiedMessages,
			}

			data, err := json.Marshal(responseData)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to serialize chat data"))
				continue
			}

			resp := internal.Response{
				Command: internal.ResponseSelectChat,
				Data:    string(data),
			}
			log.Printf("Sending data: %s", data)

			activeUsers[username] = internal.ActiveUser{
				WsConnection:   conn,
				Expiry:         time.Now().Add(2 * time.Minute),
				PublicKey:      activeUsers[username].PublicKey,
				SelectedFriend: friend,
			}

			sendJSON(conn, resp)

		case "ping":

			var dataMap map[string]string
			err := json.Unmarshal([]byte(message.Data), &dataMap)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format"))
				continue
			}

			username := dataMap["username"]

			_, exist := activeUsers[username]
			if !exist {
				conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			if time.Now().After(activeUsers[username].Expiry) {
				conn.WriteMessage(websocket.TextMessage, []byte("Session expired"))
				continue
			}

			activeUsers[username] = internal.ActiveUser{
				WsConnection:   conn,
				Expiry:         time.Now().Add(2 * time.Minute),
				PublicKey:      activeUsers[username].PublicKey,
				SelectedFriend: "",
			}

			friendList, err := queries.GetContactsByUsername(db.GetDBInstance(), username)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to get contacts"))
				continue
			}

			resp := wsresponses.ResponseSolveSuccess(activeUsers[username].PublicKey, friendList)
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

	go startUserSessionChecker()

	router := gin.Default()

	router.GET("/ws", wsHandler)

	router.Run(":8080")
}
