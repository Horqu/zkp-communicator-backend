package views

import (
	"fmt"
	"log"

	"client/internal"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gorilla/websocket"
)

var (
	privateKeyEditorResolver = new(widget.Editor)
	resolveButton            = new(widget.Clickable)
	UserPrivateKey           string
)

func LayoutResolver(gtx layout.Context, th *material.Theme, currentView *internal.AppView, wsConn *websocket.Conn, usernameLogin *string) layout.Dimensions {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			edit := material.Editor(th, privateKeyEditorResolver, "Enter private key")
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, resolveButton, "Resolve")
			if resolveButton.Clicked(gtx) {
				privateKey := privateKeyEditorResolver.Text()
				log.Printf("Wprowadzony klucz prywatny: %s\n", privateKey)
				UserPrivateKey = privateKey
				if wsConn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    fmt.Sprintf(`{"username":"%s"}`, *usernameLogin),
					}
					err := wsConn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("Failed to send solve message: %v\n", err)
					} else {
						fmt.Printf("Sent solve message for username=%s\n", *usernameLogin)
					}
				} else {
					fmt.Println("WebSocket connection is not established.")
				}
			}
			return btn.Layout(gtx)
		}),
	)

	return layout.Dimensions{}
}
