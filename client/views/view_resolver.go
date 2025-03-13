package views

import (
	"log"

	"client/internal"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	privateKeyEditorResolver = new(widget.Editor)
	resolveButton            = new(widget.Clickable)
)

func LayoutResolver(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
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
				// Tutaj można zrealizować weryfikację klucza prywatnego
			}
			return btn.Layout(gtx)
		}),
	)

	return layout.Dimensions{}
}
