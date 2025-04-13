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
	privateKeyEditorResolverFFS = new(widget.Editor)
	resolveButtonFFS            = new(widget.Clickable)
	UserPrivateKeyFFS           string
)

func LayoutResolverFFS(gtx layout.Context, th *material.Theme, currentView *internal.AppView, wsConn *websocket.Conn, usernameLogin *string) layout.Dimensions {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			edit := material.Editor(th, privateKeyEditorResolverFFS, "Enter private key")
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, resolveButtonFFS, "Resolve")
			if resolveButtonFFS.Clicked(gtx) {
				privateKey := privateKeyEditorResolverFFS.Text()
				log.Printf("Wprowadzony klucz prywatny: %s\n", privateKey)
				UserPrivateKeyFFS = privateKey
				if wsConn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    fmt.Sprintf(`{"username":"%s"}`, *usernameLogin),
					}
					err := wsConn.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send solve message: %v\n", err)
					} else {
						log.Printf("Sent solve message for username=%s\n", *usernameLogin)
					}
				} else {
					log.Println("WebSocket connection is not established.")
				}
			}
			return btn.Layout(gtx)
		}),
	)

	return layout.Dimensions{}
}
