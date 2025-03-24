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
	loginEditor        = new(widget.Editor)
	verificationOption = new(widget.Enum)
	sendButton         = new(widget.Clickable)
)

func LayoutLogin(gtx layout.Context, th *material.Theme, currentView *internal.AppView, wsConn *websocket.Conn, usernameLogin *string) layout.Dimensions {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			edit := material.Editor(th, loginEditor, "Enter login")
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rb := material.RadioButton(th, verificationOption, "methodA", "ZK-SNARK")
					return rb.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rb := material.RadioButton(th, verificationOption, "methodB", "Method B")
					return rb.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rb := material.RadioButton(th, verificationOption, "methodC", "Method C")
					return rb.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, sendButton, "Send")
			if sendButton.Clicked(gtx) {
				// Przykładowa logika wysyłania
				login := loginEditor.Text()
				method := verificationOption.Value
				log.Printf("Sending data: login=%s, method=%s\n", login, method)
				if wsConn != nil {
					msg := internal.Message{
						Command: internal.MessageLogin,
						Data:    fmt.Sprintf(`{"username":"%s","method":"%s"}`, login, method),
					}
					err := wsConn.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to send registration message: %v\n", err)
					} else {
						log.Printf("Sent login message: username=%s, publicKey=%s\n", login, method)
						*usernameLogin = login
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
