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

	ResetChan chan bool
)

func clearClientVariables() {
	usernameLogin = ""
	friendList = nil
	userPublicKey = ""
	selectedFriendPublicKey = ""
	decryptedMessagesGlobal = nil

	schnorr_E = nil
	schnorr_R = nil

	ffs_C1N = bn254.G1Affine{}
	ffs_EncryptedN = ""
	ffs_C1e = bn254.G1Affine{}
	ffs_EncryptedE = ""

	sigma_C1e = bn254.G1Affine{}
	sigma_EncryptedE = ""
	sigma_C1r = bn254.G1Affine{}
	sigma_EncryptedR = ""

}

func main() {
	go func() {

		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	ResetChan = make(chan bool)
	go internal.StartSessionTimer(ResetChan)

	app.Main()
}

func connectToWebSocket(wsURL string) error {
	if wsConn != nil {
		log.Println("WebSocket is already connected.")
		return nil
	}

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

	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			var msg internal.Response
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Not a Json message:", string(message))
				continue
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

		currentView = internal.ViewResolverSchnorr

	case internal.ResponseFFSChallenge:
		log.Println("Received challenge")
		// Rozpakowanie danych z odpowiedzi
		var responseData struct {
			C1e        string `json:"C1e"`
			EncryptedE string `json:"EncryptedE"`
			C1N        string `json:"C1N"`
			EncryptedN string `json:"EncryptedN"`
		}
		err := json.Unmarshal([]byte(msg.Data), &responseData)
		if err != nil {
			log.Printf("Failed to parse ResponseFFSChallenge data: %v", err)
			return
		}
		ffs_C1N, err = encryption.StringToPublicKey(responseData.C1N)
		if err != nil {
			log.Printf("Failed to convert C1N to G1Affine: %v", err)
			return
		}
		ffs_EncryptedN = responseData.EncryptedN
		ffs_C1e, err = encryption.StringToPublicKey(responseData.C1e)
		if err != nil {
			log.Printf("Failed to convert C1e to G1Affine: %v", err)
			return
		}
		ffs_EncryptedE = responseData.EncryptedE

		currentView = internal.ViewResolverFFS

	case internal.ResponseSigmaChallenge:
		log.Println("Received challenge")
		// Rozpakowanie danych z odpowiedzi
		var responseData struct {
			C1e        string `json:"C1e"`
			EncryptedE string `json:"EncryptedE"`
			C1r        string `json:"C1r"`
			EncryptedR string `json:"EncryptedR"`
		}
		err := json.Unmarshal([]byte(msg.Data), &responseData)
		if err != nil {
			log.Printf("Failed to parse ResponseSigmaChallenge data: %v", err)
			return
		}
		sigma_C1e, err = encryption.StringToPublicKey(responseData.C1e)
		if err != nil {
			log.Printf("Failed to convert C1e to G1Affine: %v", err)
			return
		}
		sigma_EncryptedE = responseData.EncryptedE
		sigma_C1r, err = encryption.StringToPublicKey(responseData.C1r)
		if err != nil {
			log.Printf("Failed to convert C1r to G1Affine: %v", err)
			return
		}
		sigma_EncryptedR = responseData.EncryptedR

		currentView = internal.ViewResolverSigma

	case internal.ResponseSolveSuccess:
		log.Println("Solved successfully")
		ResetChan <- true

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
		userPrivateKeyBigInt.SetString(encryption.UserPrivateKey, 10)
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

func closeWebSocket() {
	if wsConn != nil {
		err := wsConn.Close()
		if err != nil {
			log.Printf("Failed to close WebSocket connection: %v\n", err)
		} else {
			log.Println("WebSocket connection closed successfully.")
		}
		wsConn = nil
	}
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			closeWebSocket()
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			switch currentView {
			case internal.ViewMain:
				views.LayoutMain(gtx, th, &currentView)
			case internal.ViewLogin, internal.ViewRegister:
				// Połącz z WebSocket, jeśli jeszcze nie jest połączony
				if wsConn == nil {
					err := connectToWebSocket("ws://localhost:8080/ws")
					if err != nil {
						log.Printf("Failed to connect to WebSocket: %v\n", err)
					}
				}
				if currentView == internal.ViewLogin {
					views.LayoutLogin(gtx, th, &currentView, wsConn, &usernameLogin)
				} else {
					views.LayoutRegister(gtx, th, &currentView, wsConn)
				}
			case internal.ViewResolverSchnorr:
				views.LayoutResolverSchnorr(gtx, th, &currentView, wsConn, &usernameLogin, schnorr_R, schnorr_E)
			case internal.ViewResolverFFS:
				views.LayoutResolverFFS(gtx, th, &currentView, wsConn, &usernameLogin, ffs_C1N, ffs_EncryptedN, ffs_C1e, ffs_EncryptedE)
			case internal.ViewResolverSigma:
				views.LayoutResolverSigma(gtx, th, &currentView, wsConn, &usernameLogin, sigma_C1e, sigma_EncryptedE, sigma_C1r, sigma_EncryptedR)
			case internal.ViewResolver:
				views.LayoutResolver(gtx, th, &currentView, wsConn, &usernameLogin)
			case internal.ViewLoading:
				views.LayoutLoading(gtx, th, &currentView)
			case internal.ViewError:
				views.LayoutError(gtx, th, &currentView)
			case internal.ViewLogged:
				if internal.GetSessionTimeLeft() <= 0 {
					log.Println("Session expired")
					currentView = internal.ViewLogout
				}
				views.LayoutLogged(gtx, th, &currentView, wsConn, &usernameLogin, friendList, userPublicKey, selectedFriendPublicKey, decryptedMessagesGlobal, &ResetChan)
			case internal.ViewLogout:
				currentView = internal.ViewMain
				closeWebSocket()
				clearClientVariables()
			}

			e.Frame(gtx.Ops)
		}
	}
}
