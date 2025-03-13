package views

import (
	"client/internal"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	goToLoginButton    = new(widget.Clickable)
	goToRegisterButton = new(widget.Clickable)
)

func LayoutMain(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, goToLoginButton, "Login")
			btn.TextSize = unit.Sp(16)
			if goToLoginButton.Clicked(gtx) {
				*currentView = internal.ViewLogin
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, goToRegisterButton, "Register")
			btn.TextSize = unit.Sp(16)
			if goToRegisterButton.Clicked(gtx) {
				*currentView = internal.ViewRegister
			}
			return btn.Layout(gtx)
		}),
	)
}
