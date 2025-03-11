package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"

	"client/internal"
	"client/views"
)

var currentView internal.AppView = internal.ViewMain

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			switch currentView {
			case internal.ViewMain:
				// 1. Widok główny
				views.MainResolver(gtx, th, &currentView)
			case internal.ViewResolver:
				// 2. Drugi widok (klucz prywatny + przycisk "Resolve")
				views.LayoutResolver(gtx, th, &currentView)

			}

			e.Frame(gtx.Ops)
		}
	}
}
