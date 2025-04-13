package views

import (
	"fmt"
	"log"
	"math/big"

	"client/encryption"
	"client/internal"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gorilla/websocket"
)

var (
	privateKeyEditorResolverSchnorr = new(widget.Editor)
	resolveButtonSchnorr            = new(widget.Clickable)
	UserPrivateKeySchnorr           string
)

func LayoutResolverSchnorr(gtx layout.Context, th *material.Theme, currentView *internal.AppView, wsConn *websocket.Conn, usernameLogin *string, R *big.Int, E *big.Int) layout.Dimensions {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			edit := material.Editor(th, privateKeyEditorResolverSchnorr, "Enter private key")
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, resolveButtonSchnorr, "Resolve Schnorr")
			if resolveButtonSchnorr.Clicked(gtx) {
				privateKey := privateKeyEditorResolverSchnorr.Text()
				log.Printf("Wprowadzony klucz prywatny: %s\n", privateKey)
				UserPrivateKeySchnorr = privateKey
				UserPrivateKeySchnorrBigInt, err := encryption.StringToBigInt(privateKey)
				if err != nil {
					log.Printf("Failed to convert private key to big.Int: %v\n", err)
					return layout.Dimensions{}
				}
				s := encryption.GenerateSchnorrProof(UserPrivateKeySchnorrBigInt, R, E)
				if wsConn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    fmt.Sprintf(`{"username":"%s", "method":"Schnorr", "s":"%s"}`, *usernameLogin, s.String()),
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
