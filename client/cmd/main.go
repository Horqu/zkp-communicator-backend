package main

import (
	"encoding/json"
	"log"
	"math/big"
	"net/url"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/gorilla/websocket"

	"client/encryption"
	"client/internal"
	"client/views"
)

var (
	wsConn                  *websocket.Conn
	currentView             internal.AppView = internal.ViewMain
	usernameLogin           string
	friendList              []internal.SimplifiedContact
	userPublicKey           string
	selectedFriendPublicKey string
	decryptedMessagesGlobal []views.DecryptedMessage
	activeChallenge         string

	schnorr_E *big.Int
	schnorr_R *big.Int

	ffs_C1N        bn254.G1Affine
	ffs_EncryptedN string
	ffs_C1e        bn254.G1Affine
	ffs_EncryptedE string

	sigma_C1e        bn254.G1Affine
	sigma_EncryptedE string
	sigma_C1r        bn254.G1Affine
	sigma_EncryptedR string
)

func main() {
	go func() {

		if err := connectToWebSocket("ws://localhost:8080/ws"); err != nil {
			log.Println("WebSocket connection error:", err)
		}

		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func connectToWebSocket(wsURL string) error {
	u, err := url.Parse(wsURL)
	if err != nil {
		return err
	}

	log.Println("Connecting to", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	wsConn = c
	log.Println("WebSocket connected.")

	// Gorutyna do odbierania wiadomości
	go func() {
		for {
			var msg internal.Response
			err := wsConn.ReadJSON(&msg)
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			handleMessage(msg)
		}
	}()

	return nil
}

func handleMessage(msg internal.Response) {
	log.Printf("Received message: Command=%s, Data=%s\n", msg.Command, msg.Data)

	// Obsługa różnych typów wiadomości
	switch msg.Command {
	case internal.ResponseRegisterSuccess:
		log.Println("Registered successfully")
		currentView = internal.ViewMain

	case internal.ResponseSchnorrChallenge:
		log.Println("Received challenge")
		// Rozpakowanie danych z odpowiedzi
		var responseData struct {
			R *big.Int `json:"R"`
			E *big.Int `json:"E"`
		}
		err := json.Unmarshal([]byte(msg.Data), &responseData)
		if err != nil {
			log.Printf("Failed to parse ResponseSchnorrChallenge data: %v", err)
			return
		}
		schnorr_E = responseData.E
		schnorr_R = responseData.R

		activeChallenge = string(internal.ResponseSchnorrChallenge)
		currentView = internal.ViewResolverSchnorr

	case internal.ResponseFFSChallenge:
		log.Println("Received challenge")
		// Rozpakowanie danych z odpowiedzi
		var responseData struct {
			C1e        bn254.G1Affine `json:"C1e"`
			EncryptedE string         `json:"EncryptedE"`
			C1N        bn254.G1Affine `json:"C1N"`
			EncryptedN string         `json:"EncryptedN"`
		}
		err := json.Unmarshal([]byte(msg.Data), &responseData)
		if err != nil {
			log.Printf("Failed to parse ResponseFFSChallenge data: %v", err)
			return
		}
		ffs_C1N = responseData.C1N
		ffs_EncryptedN = responseData.EncryptedN
		ffs_C1e = responseData.C1e
		ffs_EncryptedE = responseData.EncryptedE
		activeChallenge = string(internal.ResponseFFSChallenge)

	case internal.ResponseSigmaChallenge:
		log.Println("Received challenge")
		// Rozpakowanie danych z odpowiedzi
		var responseData struct {
			C1e        bn254.G1Affine `json:"C1e"`
			EncryptedE string         `json:"EncryptedE"`
			C1r        bn254.G1Affine `json:"C1r"`
			EncryptedR string         `json:"EncryptedR"`
		}
		err := json.Unmarshal([]byte(msg.Data), &responseData)
		if err != nil {
			log.Printf("Failed to parse ResponseSigmaChallenge data: %v", err)
			return
		}
		sigma_C1e = responseData.C1e
		sigma_EncryptedE = responseData.EncryptedE
		sigma_C1r = responseData.C1r
		sigma_EncryptedR = responseData.EncryptedR
		activeChallenge = string(internal.ResponseSigmaChallenge)

	case internal.ResponseSolveSuccess:
		log.Println("Solved successfully")

		// Rozpakowanie kontaktów z odpowiedzi
		var responseData struct {
			PublicKey  string                       `json:"publicKey"`
			FriendList []internal.SimplifiedContact `json:"friendList"`
		}
		err := json.Unmarshal([]byte(msg.Data), &responseData)
		if err != nil {
			log.Printf("Failed to parse ResponseSolveSuccess data: %v", err)
			return
		}

		userPublicKey = responseData.PublicKey
		log.Printf("Updated userPublicKey: %s", userPublicKey)

		// Przypisanie kontaktów do globalnej zmiennej friendList
		friendList = responseData.FriendList
		log.Printf("Updated friendList: %+v", friendList)

		currentView = internal.ViewLogged
	case internal.ResponseSelectChat:
		log.Println("Selected chat successfully")

		type SimplifiedMessage struct {
			SenderUsername    string    `json:"senderUsername"`
			RecipientUsername string    `json:"recipientUsername"`
			C1                string    `json:"c1"`
			Content           string    `json:"content"`
			CreatedAt         time.Time `json:"createdAt"`
		}

		type responseData struct {
			FriendPublicKey string              `json:"friendPublicKey"`
			Messages        []SimplifiedMessage `json:"messages"`
		}

		var chatData responseData
		err := json.Unmarshal([]byte(msg.Data), &chatData)
		if err != nil {
			log.Printf("Failed to parse ResponseSelectChat data: %v", err)
			return
		}

		// Zapisanie friendPublicKey do zmiennej globalnej
		selectedFriendPublicKey = chatData.FriendPublicKey
		log.Printf("Updated selectedFriendPublicKey: %s", selectedFriendPublicKey)

		// Zapisanie wiadomości do zmiennej lokalnej
		var chatMessages []SimplifiedMessage = chatData.Messages
		log.Printf("Updated chatMessages: %+v", chatMessages)

		var decryptedMessages []views.DecryptedMessage
		userPrivateKeyBigInt := new(big.Int)
		userPrivateKeyBigInt.SetString(views.UserPrivateKey, 10)
		log.Printf("User private key: %s", userPrivateKeyBigInt)
		for _, message := range chatMessages {
			C1G1Affine, _ := encryption.StringToPublicKey(message.C1)

			decryptedContent := encryption.DecryptText(C1G1Affine, message.Content, userPrivateKeyBigInt)
			decryptedMessages = append(decryptedMessages, views.DecryptedMessage{
				SenderUsername:    message.SenderUsername,
				ReceipentUsername: message.RecipientUsername,
				Content:           string(decryptedContent),
				CreatedAt:         message.CreatedAt,
			})
		}
		log.Printf("Decrypted messages: %+v", decryptedMessages)
		decryptedMessagesGlobal = decryptedMessages

	case internal.ResponseCommand("example_command"):
		// Obsłuż wiadomość o komendzie "example_command"
		log.Println("Handling example_command:", msg.Data)
	case internal.ResponseCommand("another_command"):
		// Obsłuż inną komendę
		log.Println("Handling another_command:", msg.Data)
	default:
		log.Println("Unknown command:", msg.Command)
	}
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			switch currentView {
			case internal.ViewMain:
				views.LayoutMain(gtx, th, &currentView)
			case internal.ViewLogin:
				views.LayoutLogin(gtx, th, &currentView, wsConn, &usernameLogin)
			case internal.ViewRegister:
				views.LayoutRegister(gtx, th, &currentView, wsConn)
			case internal.ViewResolverSchnorr:
				views.LayoutResolverSchnorr(gtx, th, &currentView, wsConn, &usernameLogin, schnorr_R, schnorr_E)
			case internal.ViewResolver:
				views.LayoutResolver(gtx, th, &currentView, wsConn, &usernameLogin)
			case internal.ViewLoading:
				views.LayoutLoading(gtx, th, &currentView)
			case internal.ViewError:
				views.LayoutError(gtx, th, &currentView)
			case internal.ViewLogged:
				views.LayoutLogged(gtx, th, &currentView, wsConn, &usernameLogin, friendList, userPublicKey, selectedFriendPublicKey, decryptedMessagesGlobal)
			}

			e.Frame(gtx.Ops)
		}
	}
}
