package views

import (
	"log"

	"client/internal"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	loginEditor        = new(widget.Editor)
	verificationOption = new(widget.Enum)
	sendButton         = new(widget.Clickable)
)

func LayoutLogin(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			edit := material.Editor(th, loginEditor, "Enter login")
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rb := material.RadioButton(th, verificationOption, "methodA", "Method A")
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
				log.Printf("Wysyłam dane: login=%s, method=%s\n", login, method)
				*currentView = internal.ViewResolver
			}
			return btn.Layout(gtx)
		}),
	)
	return layout.Dimensions{}
}
