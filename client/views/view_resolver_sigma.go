package views

import (
	"fmt"
	"log"

	"client/encryption"
	"client/internal"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/gorilla/websocket"
)

var (
	privateKeyEditorResolverSigma = new(widget.Editor)
	resolveButtonSigma            = new(widget.Clickable)
	UserPrivateKeySigma           string
)

func LayoutResolverSigma(gtx layout.Context, th *material.Theme, currentView *internal.AppView, wsConn *websocket.Conn, usernameLogin *string, sigma_C1e bn254.G1Affine, sigma_EncryptedE string, sigma_C1r bn254.G1Affine, sigma_EncryptedR string) layout.Dimensions {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			edit := material.Editor(th, privateKeyEditorResolverSigma, "Enter private key")
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, resolveButtonSigma, "Resolve")
			if resolveButtonSigma.Clicked(gtx) {
				privateKey := privateKeyEditorResolverSigma.Text()
				log.Printf("Wprowadzony klucz prywatny: %s\n", privateKey)
				encryption.UserPrivateKey = privateKey
				UserPrivateKeySigma = privateKey

				UserPrivateKeySigmaBigInt, err := encryption.StringToBigInt(privateKey)
				if err != nil {
					log.Printf("Failed to convert private key to big.Int: %v\n", err)
					return layout.Dimensions{}
				}

				decryptedE := encryption.DecryptText(sigma_C1e, sigma_EncryptedE, UserPrivateKeySigmaBigInt)
				decryptedEBigInt, err := encryption.StringToBigInt(decryptedE)
				if err != nil {
					log.Printf("Failed to convert decrypted E to big.Int: %v\n", err)
					return layout.Dimensions{}
				}

				decryptedR := encryption.DecryptText(sigma_C1r, sigma_EncryptedR, UserPrivateKeySigmaBigInt)
				decryptedRBigInt, err := encryption.StringToBigInt(decryptedR)
				if err != nil {
					log.Printf("Failed to convert decrypted R to big.Int: %v\n", err)
					return layout.Dimensions{}
				}

				log.Printf("Decrypted E: %s\n", decryptedEBigInt)
				log.Printf("Decrypted R: %s\n", decryptedRBigInt)
				s := encryption.GenerateSigmaProof(UserPrivateKeySigmaBigInt, decryptedEBigInt, decryptedRBigInt)
				privateKeyEditorResolverSigma.SetText("") // Reset the private key editor after sending
				if wsConn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    fmt.Sprintf(`{"username":"%s", "method":"Sigma", "s":"%s"}`, *usernameLogin, s.String()),
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
