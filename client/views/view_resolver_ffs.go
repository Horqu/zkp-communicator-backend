package views

import (
	"encoding/json"
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
	privateKeyEditorResolverFFS = new(widget.Editor)
	resolveButtonFFS            = new(widget.Clickable)
	UserPrivateKeyFFS           string
)

func LayoutResolverFFS(gtx layout.Context, th *material.Theme, currentView *internal.AppView, wsConn *websocket.Conn, usernameLogin *string, ffs_C1N bn254.G1Affine, ffs_EncryptedN string, ffs_C1e bn254.G1Affine, ffs_EncryptedE string) layout.Dimensions {
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
				encryption.UserPrivateKey = privateKey
				UserPrivateKeyFFS = privateKey
				UserPrivateKeyFFSBigInt, err := encryption.StringToBigInt(privateKey)
				if err != nil {
					log.Printf("Failed to convert private key to big.Int: %v\n", err)
					return layout.Dimensions{}
				}
				decryptedN := encryption.DecryptText(ffs_C1N, ffs_EncryptedN, UserPrivateKeyFFSBigInt)
				decryptedNBigInt, err := encryption.StringToBigInt(decryptedN)
				if err != nil {
					log.Printf("Failed to convert decrypted N to big.Int: %v\n", err)
					return layout.Dimensions{}
				}
				decryptedE := encryption.DecryptText(ffs_C1e, ffs_EncryptedE, UserPrivateKeyFFSBigInt)
				decryptedEBigInt, err := encryption.StringToBigInt(decryptedE)
				if err != nil {
					log.Printf("Failed to convert decrypted E to big.Int: %v\n", err)
					return layout.Dimensions{}
				}

				log.Printf("Decrypted N: %s\n", decryptedN)
				log.Printf("Decrypted E: %s\n", decryptedEBigInt.Text(2))

				xList, yList, v := encryption.GenerateFeigeFiatShamirProof(UserPrivateKeyFFSBigInt, decryptedNBigInt, decryptedEBigInt)

				xListJson, err := internal.BigIntSliceToJSONString(xList)
				if err != nil {
					log.Printf("Failed to convert xList to JSON: %v\n", err)
					return layout.Dimensions{}
				}
				yListJson, err := internal.BigIntSliceToJSONString(yList)
				if err != nil {
					log.Printf("Failed to convert yList to JSON: %v\n", err)
					return layout.Dimensions{}
				}
				// Utwórz mapę danych
				dataMap := map[string]string{
					"username": *usernameLogin,
					"method":   "FeigeFiatShamir",
					"xList":    xListJson,
					"yList":    yListJson,
					"v":        v.String(),
				}

				// Serializuj mapę do JSON
				jsonBytes, err := json.Marshal(dataMap)
				if err != nil {
					log.Printf("Failed to serialize data map to JSON: %v\n", err)
					return layout.Dimensions{}
				}

				if wsConn != nil {
					msg := internal.Message{
						Command: internal.MessageSolve,
						Data:    string(jsonBytes),
					}
					log.Printf("JSON size: %d\n", len(msg.Data))
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
