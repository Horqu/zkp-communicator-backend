package views

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"client/internal"
)

// Globalny przycisk do obsługi powrotu
var backButton = new(widget.Clickable)

func LayoutError(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Nagłówek błędu
			title := material.Label(th, unit.Sp(20), "Error occurred!")
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk "Back to main page"
			btn := material.Button(th, backButton, "Back to main page")
			if backButton.Clicked(gtx) {
				*currentView = internal.ViewMain
			}
			return btn.Layout(gtx)
		}),
	)
}
