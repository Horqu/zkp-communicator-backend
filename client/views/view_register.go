package views

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"client/internal"
)

// Przycisk do wygenerowania kluczy
var generateKeysButton = new(widget.Clickable)

// Edytory (zamiast labeli) do wyświetlania kluczy
var (
	privateKeyEditor = new(widget.Editor)
	publicKeyEditor  = new(widget.Editor)
)

// Formularz użytkownika
var (
	usernameEditor = new(widget.Editor)
	pubKeyEditor   = new(widget.Editor)
	registerButton = new(widget.Clickable)
)

// LayoutRegister wyświetla przycisk do generowania kluczy, pokazuje klucze prywatny/publiczny i formularz rejestracji
func LayoutRegister(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk "Generate keys"
			btn := material.Button(th, generateKeysButton, "Generate keys")
			if generateKeysButton.Clicked(gtx) {
				privateKeyEditor.SetText("PRIVATE_KEY_EXAMPLE")
				publicKeyEditor.SetText("PUBLIC_KEY_EXAMPLE")
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Wyświetlamy klucze tylko jeśli istnieją w edytorach
			if privateKeyEditor.Text() != "" && publicKeyEditor.Text() != "" {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(14), "Private key:")
						lbl.Alignment = text.Start
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						ed := material.Editor(th, privateKeyEditor, "")
						// Nie zmieniamy textu w kodzie, więc w praktyce jest to „read-only”
						return ed.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(14), "Public key:")
						lbl.Alignment = text.Start
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						ed := material.Editor(th, publicKeyEditor, "")
						return ed.Layout(gtx)
					}),
				)
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Formularz: username, public key
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), "Username:")
					lbl.Alignment = text.Start
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(th, usernameEditor, "Enter username")
					return ed.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), "Public key:")
					lbl.Alignment = text.Start
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					pubEd := material.Editor(th, pubKeyEditor, "Paste or type your public key")
					return pubEd.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk "Register"
			btn := material.Button(th, registerButton, "Register")
			if registerButton.Clicked(gtx) {
				u := usernameEditor.Text()
				pk := pubKeyEditor.Text()
				fmt.Printf("Register user=%s with publicKey=%s\n", u, pk)
			}
			return btn.Layout(gtx)
		}),
	)
}
