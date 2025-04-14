package views

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gorilla/websocket"

	"client/encryption"
	"client/internal"
)

type DecryptedMessage struct {
	SenderUsername    string
	ReceipentUsername string
	Content           string
	CreatedAt         time.Time
}

// Zmienne do logiki ekranu Logged
var (
	// Górny pasek
	refreshButton = new(widget.Clickable)
	logoutButton  = new(widget.Clickable)

	// Lista znajomych i zaznaczenie aktualnego
	friendList     = []string{}
	selectedFriend = -1

	// Dolny pasek: dodawanie znajomego i wysyłanie wiadomości
	newFriendEditor   = new(widget.Editor)
	addFriendButton   = new(widget.Clickable)
	messageEditor     = new(widget.Editor)
	sendMessageButton = new(widget.Clickable)

	friendButtons []*widget.Clickable
	chatList      widget.List

	// WebSocket connection
	wsConnGlobal *websocket.Conn

	// Login
	usernameLoginGlobal string

	decryptedMessages       []DecryptedMessage
	selectedFriendPublicKey string
	userPublicKeyGlobal     string
)

func clearLoggedVariables() {
	friendList = []string{}
	selectedFriend = -1
	newFriendEditor.SetText("")
	messageEditor.SetText("")
	friendButtons = nil
	chatList = widget.List{}
	wsConnGlobal = nil
	usernameLoginGlobal = ""
	selectedFriendPublicKey = ""
	userPublicKeyGlobal = ""
}

// LayoutLogged - główny ekran (chat) po zalogowaniu
func LayoutLogged(gtx layout.Context, th *material.Theme, currentView *internal.AppView, wsConn *websocket.Conn, usernameLogin *string, contactList []internal.SimplifiedContact, userPublicKey string, friendPublicKey string, messages []DecryptedMessage, resetChan *chan bool) layout.Dimensions {
	userPublicKeyGlobal = userPublicKey
	selectedFriendPublicKey = friendPublicKey
	decryptedMessages = messages
	friendList = make([]string, len(contactList))
	wsConnGlobal = wsConn
	usernameLoginGlobal = *usernameLogin
	for i, contact := range contactList {
		friendList[i] = contact.Username
	}

	// Układ pionowy (top bar, środek, bottom bar)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Górny pasek
			return layoutTopBar(gtx, th, currentView, resetChan)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Środek aplikacji
			return layoutMiddle(gtx, th, resetChan)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Dolny pasek
			return layoutBottomBar(gtx, th, resetChan)
		}),
	)
}

// layoutTopBar tworzy pasek z "session time left", przyciskiem Refresh i Logout
func layoutTopBar(gtx layout.Context, th *material.Theme, currentView *internal.AppView, resetChan *chan bool) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sessionTimeLeft := internal.GetSessionTimeLeft()
			minutes := sessionTimeLeft / 60
			seconds := sessionTimeLeft % 60
			sessionLabel := material.Label(th, unit.Sp(20), fmt.Sprintf("Session time left: %02d:%02d", minutes, seconds))
			sessionLabel.Alignment = text.Middle
			sessionLabel.Color = color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF} // Czerwony kolor dla lepszej widoczności
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
			}.Layout(gtx, sessionLabel.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk Refresh
			btn := material.Button(th, refreshButton, "Refresh")
			if refreshButton.Clicked(gtx) {
				*resetChan <- true
				// Wysyłanie wiadomości MessageRefresh do serwera
				if wsConnGlobal != nil {
					msg := internal.Message{
						Command: internal.MessageRefresh,
						Data:    fmt.Sprintf(`{"username":"%s"}`, usernameLoginGlobal),
					}
					err := wsConnGlobal.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send refresh message: %v\n", err)
					} else {
						log.Printf("Sent refresh message for user: %s\n", usernameLoginGlobal)
					}
				} else {
					log.Println("WebSocket connection is not established.")
				}
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk Logout
			btn := material.Button(th, logoutButton, "Logout")
			if logoutButton.Clicked(gtx) {
				*resetChan <- true

				// Rozłącz WebSocket
				if wsConnGlobal != nil {
					err := wsConnGlobal.Close()
					if err != nil {
						log.Printf("Failed to close WebSocket connection: %v\n", err)
					} else {
						log.Println("WebSocket connection closed successfully.")
					}
					wsConnGlobal = nil // Wyzeruj wsConnGlobal, aby wskazywał na brak połączenia
				}
				clearLoggedVariables()
				// Powrót do ekranu logowania
				*currentView = internal.ViewLogout

			}
			return btn.Layout(gtx)
		}),
	)
}

// layoutMiddle - środkowa część aplikacji: lista znajomych (lewa) i chat z wybranym znajomym (prawa)
func layoutMiddle(gtx layout.Context, th *material.Theme, resetChan *chan bool) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Lista znajomych, pionowo
			return layoutFriendsList(gtx, th, resetChan)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Chat z aktualnie zaznaczonym znajomym
			return layoutChat(gtx, th)
		}),
	)
}

func layoutFriendsList(gtx layout.Context, th *material.Theme, resetChan *chan bool) layout.Dimensions {
	// Jeśli friendButtons jest puste lub długość się nie zgadza, inicjalizujemy ponownie
	if len(friendButtons) != len(friendList) {
		friendButtons = make([]*widget.Clickable, len(friendList))
		for i := range friendButtons {
			friendButtons[i] = new(widget.Clickable)
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		friendItems(gtx, th, resetChan)...,
	)
}

func friendItems(gtx layout.Context, th *material.Theme, resetChan *chan bool) []layout.FlexChild {
	var children []layout.FlexChild
	for i, friend := range friendList {
		index := i
		button := friendButtons[i] // Pobieramy już istniejący widget.Clickable

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			item := material.Button(th, button, friend)
			if button.Clicked(gtx) {
				*resetChan <- true
				log.Printf("Selected friend: %s\n", friend)
				selectedFriend = index
				friendName := friendList[selectedFriend]
				log.Printf("Selected friend: %s\n", friendName)
				if wsConnGlobal != nil {
					msg := internal.Message{
						Command: internal.MessageSelectChat,
						Data:    fmt.Sprintf(`{"username":"%s","friend":"%s"}`, usernameLoginGlobal, friendName),
					}
					err := wsConnGlobal.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send select chat message: %v\n", err)
					} else {
						log.Printf("Sent select chat message: username=%s, friend=%s\n", usernameLoginGlobal, friendName)
					}
				} else {
					log.Println("WebSocket connection is not established.")
				}
			}

			// Podświetlenie wybranego znajomego
			if index == selectedFriend {
				item.Background = color.NRGBA{R: 0x00, G: 0xAA, B: 0xFF, A: 0xFF}
				item.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
			} else {
				item.Background = color.NRGBA{R: 0xBB, G: 0xBB, B: 0xBB, A: 0xFF}
				item.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
			}

			item.Inset = layout.UniformInset(unit.Dp(4))
			return item.Layout(gtx)
		}))
	}
	return children
}

func layoutChat(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if selectedFriend < 0 || selectedFriend >= len(friendList) {
		// Jeśli nic nie wybrano, pokaż placeholder
		lbl := material.Label(th, unit.Sp(16), "Select a friend to chat")
		lbl.Alignment = text.Middle
		return lbl.Layout(gtx)
	}

	// Lista wiadomości
	friendName := friendList[selectedFriend]
	chatLabel := material.Label(th, unit.Sp(16), "Chat with "+friendName)
	chatLabel.Alignment = text.Start

	// Używamy widget.List do obsługi przewijania
	chatList.Axis = layout.Vertical

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Nagłówek chatu
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return chatLabel.Layout(gtx)
		}),
		// Lista wiadomości z przewijaniem
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return chatList.Layout(gtx, len(decryptedMessages), func(gtx layout.Context, i int) layout.Dimensions {
				// Pobieramy wiadomość
				message := decryptedMessages[i]

				// Wyświetlamy wiadomość
				var bubbleColor color.NRGBA
				var textColor color.NRGBA
				if message.SenderUsername == usernameLoginGlobal {
					// Wiadomość wysłana przez użytkownika
					bubbleColor = color.NRGBA{R: 0xD1, G: 0xF7, B: 0xC4, A: 0xFF} // Zielony
					textColor = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}   // Czarny
				} else {
					// Wiadomość odebrana
					bubbleColor = color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF} // Szary
					textColor = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}   // Czarny
				}

				// Styl wiadomości
				bubble := material.Label(th, unit.Sp(14), fmt.Sprintf("%s: %s", message.SenderUsername, message.Content))
				bubble.Color = textColor

				return layout.Inset{
					Top:    unit.Dp(4),
					Bottom: unit.Dp(4),
					Left:   unit.Dp(8),
					Right:  unit.Dp(8),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							// Tło wiadomości
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return widget.Border{
									Color:        bubbleColor,
									CornerRadius: unit.Dp(4),
									Width:        unit.Dp(1),
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return bubble.Layout(gtx)
								})
							})
						}),
					)
				})
			})
		}),
	)
}

// layoutBottomBar - pole do wpisania nowego znajomego i wiadomości
func layoutBottomBar(gtx layout.Context, th *material.Theme, resetChan *chan bool) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Pole do wpisania nowego znajomego
			edit := material.Editor(th, newFriendEditor, "Add new friend")
			// edit.SingleLine = true
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk "Add friend"
			btn := material.Button(th, addFriendButton, "Add")
			if addFriendButton.Clicked(gtx) {
				*resetChan <- true
				// Dodaj znajomego do listy
				name := newFriendEditor.Text()
				if name != "" {
					// friendList = append(friendList, name)
					log.Printf("Add friend: %s\n", name)

					if wsConnGlobal != nil {
						msg := internal.Message{
							Command: internal.MessageAddFriend,
							Data:    fmt.Sprintf(`{"username":"%s","friend":"%s"}`, usernameLoginGlobal, name),
						}
						err := wsConnGlobal.WriteJSON(msg)
						if err != nil {
							log.Printf("Failed to send add friend message: %v\n", err)
						} else {
							log.Printf("Added friend %s for user %s\n", name, usernameLoginGlobal)
						}
					} else {
						log.Println("WebSocket connection is not established.")
					}

					newFriendEditor.SetText("")
				}
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Pole do wpisania wiadomości
			edit := material.Editor(th, messageEditor, "Type message")
			// edit.SingleLine = true
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk "Send" do aktualnie wybranego znajomego
			btn := material.Button(th, sendMessageButton, "Send")
			if sendMessageButton.Clicked(gtx) {
				*resetChan <- true
				if selectedFriend >= 0 && selectedFriend < len(friendList) {
					msgText := messageEditor.Text()
					friendUsername := friendList[selectedFriend]

					G1AffineUserPublicKey, _ := encryption.StringToPublicKey(userPublicKeyGlobal)
					C1ForSender, ContentForSender := encryption.EncryptText(msgText, &G1AffineUserPublicKey)
					C1ForSenderString := encryption.PublicKeyToString(C1ForSender)

					G1AffineFriendPublicKey, _ := encryption.StringToPublicKey(selectedFriendPublicKey)
					C1ForFriend, ContentForFriend := encryption.EncryptText(msgText, &G1AffineFriendPublicKey)
					C1ForFriendString := encryption.PublicKeyToString(C1ForFriend)

					log.Printf("Send message from %s to %s\n", usernameLoginGlobal, friendUsername)
					if wsConnGlobal != nil {
						msg := internal.Message{
							Command: internal.MessageSendMessage,
							Data:    fmt.Sprintf(`{"username":"%s","friend":"%s","c1user":"%s","contentuser":"%s","c1friend":"%s","contentfriend":"%s"}`, usernameLoginGlobal, friendUsername, C1ForSenderString, ContentForSender, C1ForFriendString, ContentForFriend),
						}
						err := wsConnGlobal.WriteJSON(msg)
						if err != nil {
							log.Printf("Failed to send select chat message: %v\n", err)
						} else {
							log.Printf("Sent chat message from username=%s to friend=%s\n", usernameLoginGlobal, friendUsername)
						}
					} else {
						log.Println("WebSocket connection is not established.")
					}
					messageEditor.SetText("")
				}
			}
			return btn.Layout(gtx)
		}),
	)
}
