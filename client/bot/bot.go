package main

import (
	botutils "client/bot-utils"
	"client/encryption"
	"client/internal"
	"client/views"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/gorilla/websocket"
)

var LOGIN_ONLY = false

func registerFlow() {
	serverURL := "ws://192.168.0.130:8080/ws"

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to WebSocket server")

	// Kanał do sygnalizowania sukcesu rejestracji
	registerSuccessChan := make(chan bool)

	// Gorutyna do odbierania wiadomości
	go func() {
		defer close(registerSuccessChan)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}

			var msg internal.Response
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			if msg.Command == internal.ResponseRegisterSuccess {
				log.Println("Registration successful")
				registerSuccessChan <- true // Sygnalizuj sukces
				return
			}
		}
	}()

	// Generowanie danych bota
	botUsername := botutils.GenerateBotUsername()
	privateKey := encryption.GeneratePrivateKey()
	publicKey := encryption.GeneratePublicKey(privateKey)
	publicKeyString := encryption.PublicKeyToString(publicKey)
	botutils.SaveUsernameAndPrivateKey(botUsername, privateKey.String())

	// Wysłanie wiadomości rejestracyjnej
	if conn != nil {
		msg := internal.Message{
			Command: internal.MessageRegister,
			Data:    fmt.Sprintf(`{"username":"%s","publicKey":"%s"}`, botUsername, publicKeyString),
		}
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Failed to send registration message: %v\n", err)
		} else {
			log.Printf("Sent registration message: username=%s, publicKey=%s\n", botUsername, publicKeyString)
		}
	} else {
		log.Println("WebSocket connection is not established.")
		return
	}

	// Czekaj na odpowiedź lub timeout
	select {
	case <-registerSuccessChan:
		log.Println("Bot successfully registered and exiting.")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for registration response. Exiting.")
	}
}

func loginFlow(loginMethod string, disconnectTime string) {

	disconnectMinutes, err := strconv.Atoi(disconnectTime)
	if err != nil {
		log.Fatalf("Invalid disconnectTime value: %v", err)
	}
	go func() {
		<-time.After(time.Duration(disconnectMinutes) * time.Minute)
		log.Printf("%s minutes have passed. Exiting bot.\n", disconnectTime)
		os.Exit(0) // Wymuś zakończenie programu
	}()

	var loggedIn bool
	loginMethod = botutils.GetLoginMethod(loginMethod)

	var username string
	var usernamePrivateKey string
	var userPublicKey string

	var friendList []internal.SimplifiedContact

	var selectedFriendUsername string
	var selectedFriendPublicKey string

	var schnorr_E *big.Int
	var schnorr_R *big.Int

	var ffs_C1N bn254.G1Affine
	var ffs_EncryptedN string
	var ffs_C1e bn254.G1Affine
	var ffs_EncryptedE string

	var sigma_C1e bn254.G1Affine
	var sigma_EncryptedE string
	var sigma_C1r bn254.G1Affine
	var sigma_EncryptedR string

	filePath := "bot_credentials.json"
	credential, err := botutils.LoadRandomBotCredential(filePath)
	if err != nil {
		log.Fatalf("Failed to load random bot credential: %v", err)
	}
	username = credential.Username
	usernamePrivateKey = credential.PrivateKey

	serverURL := "ws://192.168.0.130:8080/ws"

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to WebSocket server")

	connectionClosedChan := make(chan bool)

	// Gorutyna do odbierania wiadomości
	go func() {
		defer close(connectionClosedChan)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				connectionClosedChan <- true
				return
			}

			var msg internal.Response
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			switch msg.Command {
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

				userPrivateKeyBigInt, err := encryption.StringToBigInt(usernamePrivateKey)
				if err != nil {
					log.Printf("Failed to convert private key to big.Int: %v", err)
					return
				}
				s := encryption.GenerateSchnorrProof(userPrivateKeyBigInt, schnorr_E, schnorr_R)
				if conn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    fmt.Sprintf(`{"username":"%s", "method":"Schnorr", "s":"%s"}`, username, s.String()),
					}
					err := conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send solve message: %v\n", err)
					} else {
						log.Printf("Sent solve message for username=%s\n", username)
					}
				} else {
					log.Println("WebSocket connection is not established.")
				}

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

				userPrivateKeyBigInt, err := encryption.StringToBigInt(usernamePrivateKey)
				if err != nil {
					log.Printf("Failed to convert private key to big.Int: %v", err)
					return
				}

				decryptedN := encryption.DecryptText(ffs_C1N, ffs_EncryptedN, userPrivateKeyBigInt)
				decryptedNBigInt, err := encryption.StringToBigInt(decryptedN)
				if err != nil {
					log.Printf("Failed to convert decrypted N to big.Int: %v\n", err)
					return
				}
				decryptedE := encryption.DecryptText(ffs_C1e, ffs_EncryptedE, userPrivateKeyBigInt)
				decryptedEBigInt, err := encryption.StringToBigInt(decryptedE)
				if err != nil {
					log.Printf("Failed to convert decrypted E to big.Int: %v\n", err)
					return
				}

				xList, yList, v := encryption.GenerateFeigeFiatShamirProof(userPrivateKeyBigInt, decryptedNBigInt, decryptedEBigInt)

				xListJson, err := internal.BigIntSliceToJSONString(xList)
				if err != nil {
					log.Printf("Failed to convert xList to JSON: %v\n", err)
					return
				}

				yListJson, err := internal.BigIntSliceToJSONString(yList)
				if err != nil {
					log.Printf("Failed to convert yList to JSON: %v\n", err)
					return
				}

				dataMap := map[string]string{
					"username": username,
					"method":   "FeigeFiatShamir",
					"xList":    xListJson,
					"yList":    yListJson,
					"v":        v.String(),
				}

				// Serializuj mapę do JSON
				jsonBytes, err := json.Marshal(dataMap)
				if err != nil {
					log.Printf("Failed to serialize data map to JSON: %v\n", err)
					return
				}

				if conn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    string(jsonBytes),
					}
					err := conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send solve message: %v\n", err)
					} else {
						log.Printf("Sent solve message for username=%s\n", username)
					}
				} else {
					log.Println("WebSocket connection is not established.")
				}

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

				userPrivateKeyBigInt, err := encryption.StringToBigInt(usernamePrivateKey)
				if err != nil {
					log.Printf("Failed to convert private key to big.Int: %v", err)
					return
				}

				decryptedE := encryption.DecryptText(sigma_C1e, sigma_EncryptedE, userPrivateKeyBigInt)
				decryptedEBigInt, err := encryption.StringToBigInt(decryptedE)
				if err != nil {
					log.Printf("Failed to convert decrypted E to big.Int: %v\n", err)
					return
				}

				decryptedR := encryption.DecryptText(sigma_C1r, sigma_EncryptedR, userPrivateKeyBigInt)
				decryptedRBigInt, err := encryption.StringToBigInt(decryptedR)
				if err != nil {
					log.Printf("Failed to convert decrypted R to big.Int: %v\n", err)
					return
				}

				s := encryption.GenerateSigmaProof(userPrivateKeyBigInt, decryptedEBigInt, decryptedRBigInt)
				if conn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    fmt.Sprintf(`{"username":"%s", "method":"Sigma", "s":"%s"}`, username, s.String()),
					}
					err := conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send solve message: %v\n", err)
					} else {
						log.Printf("Sent solve message for username=%s\n", username)
					}
				} else {
					log.Println("WebSocket connection is not established.")
				}

			case internal.ResponseSolveSuccess:
				if LOGIN_ONLY {
					log.Println("Login successful, exiting.")
					if conn != nil {
						err := conn.Close()
						if err != nil {
							log.Printf("Failed to close WebSocket connection: %v\n", err)
						} else {
							log.Println("WebSocket connection closed successfully.")
						}
						conn = nil
					}
					os.Exit(0)
				}

				if !loggedIn {
					log.Println("Login successful")
					loggedIn = true
				}

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
				friendList = responseData.FriendList

			case internal.ResponseSelectChat:

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

				selectedFriendPublicKey = chatData.FriendPublicKey

				var chatMessages []SimplifiedMessage = chatData.Messages

				var decryptedMessages []views.DecryptedMessage
				userPrivateKeyBigInt := new(big.Int)
				userPrivateKeyBigInt.SetString(usernamePrivateKey, 10)

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

			}

		}
	}()

	if conn != nil {
		msg := internal.Message{
			Command: internal.MessageLogin,
			Data:    fmt.Sprintf(`{"username":"%s","method":"%s"}`, username, loginMethod),
		}
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Failed to send registration message: %v\n", err)
		} else {
			log.Printf("Sent login message: username=%s, method=%s\n", username, loginMethod)
		}
	} else {
		log.Println("WebSocket connection is not established.")
	}

	for {
		select {
		case <-connectionClosedChan:
			log.Println("Connection closed. Exiting.")
			return
		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
			if loggedIn {
				time.Sleep(time.Duration(2+rand.Intn(8)) * time.Second) // 5-20s to perform action | 2-10s
				randomValue := rand.Intn(100)

				switch {
				case randomValue < 15: // 15% szansa na "refresh"
					log.Println("Performing action: Refresh")
					msg := internal.Message{
						Command: internal.MessageRefresh,
						Data:    fmt.Sprintf(`{"username":"%s"}`, username),
					}
					err := conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send refresh message: %v\n", err)
					}

				case randomValue < 25: // 10% szansa na "add friend"
					log.Println("Performing action: Add Friend")
					friendCredential, err := botutils.LoadRandomBotCredential(filePath)
					if err != nil {
						log.Fatalf("Failed to load random bot credential: %v", err)
					}
					friendToAdd := friendCredential.Username

					if !botutils.IsFriendAlreadyAdded(friendList, friendToAdd) && friendToAdd != username {
						msg := internal.Message{
							Command: internal.MessageAddFriend,
							Data:    fmt.Sprintf(`{"username":"%s", "friend":"%s"}`, username, friendToAdd),
						}
						err := conn.WriteJSON(msg)
						if err != nil {
							log.Printf("Failed to send add friend message: %v\n", err)
						}
					} else {
						log.Println("No friends available to add.")
					}

				case randomValue < 75: // 50% szansa na "send message to friend"
					log.Println("Performing action: Send Message to Friend")
					if selectedFriendUsername != "" {

						msgText := botutils.GenerateBotMessage()

						G1AffineUserPublicKey, _ := encryption.StringToPublicKey(userPublicKey)
						C1ForSender, ContentForSender := encryption.EncryptText(msgText, &G1AffineUserPublicKey)
						C1ForSenderString := encryption.PublicKeyToString(C1ForSender)

						G1AffineFriendPublicKey, _ := encryption.StringToPublicKey(selectedFriendPublicKey)
						C1ForFriend, ContentForFriend := encryption.EncryptText(msgText, &G1AffineFriendPublicKey)
						C1ForFriendString := encryption.PublicKeyToString(C1ForFriend)

						if conn != nil {
							msg := internal.Message{
								Command: internal.MessageSendMessage,
								Data:    fmt.Sprintf(`{"username":"%s","friend":"%s","c1user":"%s","contentuser":"%s","c1friend":"%s","contentfriend":"%s"}`, username, selectedFriendUsername, C1ForSenderString, ContentForSender, C1ForFriendString, ContentForFriend),
							}
							err := conn.WriteJSON(msg)
							if err != nil {
								log.Printf("Failed to send select chat message: %v\n", err)
							} else {
								log.Printf("Sent chat message from username=%s to friend=%s\n", username, selectedFriendUsername)
							}
						} else {
							log.Println("WebSocket connection is not established.")
						}

					} else {
						log.Println("No friend selected to send a message.")
					}

				case randomValue < 99: // 24% szansa na "select friend"
					log.Println("Performing action: Select Friend")
					if len(friendList) > 0 {
						friend := friendList[rand.Intn(len(friendList))]
						selectedFriendUsername = friend.Username
						msg := internal.Message{
							Command: internal.MessageSelectChat,
							Data:    fmt.Sprintf(`{"username":"%s", "friendUsername":"%s"}`, username, selectedFriendUsername),
						}
						err := conn.WriteJSON(msg)
						if err != nil {
							log.Printf("Failed to select friend: %v\n", err)
						}
					} else {
						log.Println("No friends available to select.")
					}

				default: // 1% szansa na "logout"
					log.Println("Performing action: Logout")
					if conn != nil {
						err := conn.Close()
						if err != nil {
							log.Printf("Failed to close WebSocket connection: %v\n", err)
						} else {
							log.Println("WebSocket connection closed successfully.")
						}
						conn = nil
					}
				}
			}
		}

	}
}

func main() {

	var loginMethod string
	var disconnectTime string

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run bot.go <login|register> <loginMethod> <disconnectTime> <SetLoginOnly>")
	}

	if len(os.Args) > 2 {
		loginMethod = os.Args[2]
	} else {
		loginMethod = "random"
	}

	if len(os.Args) > 3 {
		disconnectTime = os.Args[3]
	} else {
		disconnectTime = "60"
	}

	if len(os.Args) > 4 {
		LOGIN_ONLY = true
	}

	action := os.Args[1]
	if action != "login" && action != "register" {
		log.Fatal("Invalid action. Use 'login' or 'register'.")
	}
	if action == "login" {
		log.Printf("Login action selected. Method: %s, DisconnectTime: %s", loginMethod, disconnectTime)
		loginFlow(loginMethod, disconnectTime)
		return
	}
	log.Println("Register action selected.")
	registerFlow()
}
