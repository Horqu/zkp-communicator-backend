package views

import (
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"client/internal"
)

// Zmienne do wyświetlania kropek
var (
	lastFrame time.Time
	accTime   float32
	dotsCount int
)

// LayoutLoading wyświetla napis "Loading" z rosnącą liczbą kropek (1..3) co sekundę
func LayoutLoading(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
	now := time.Now()
	if !lastFrame.IsZero() {
		dt := float32(now.Sub(lastFrame).Seconds())
		accTime += dt
		// Co 1 sekundę zwiększamy liczbę kropek (modulo 3)
		if accTime >= 1.0 {
			dotsCount++
			if dotsCount > 3 {
				dotsCount = 1
			}
			accTime = 0
		}
	}
	lastFrame = now

	labelText := "Loading" + strings.Repeat(".", dotsCount)

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(20), labelText)
		lbl.Alignment = text.Middle
		return lbl.Layout(gtx)
	})
}
