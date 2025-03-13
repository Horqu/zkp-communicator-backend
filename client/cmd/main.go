package main

import (
	"log"
	"net/url"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/gorilla/websocket"

	"client/internal"
	"client/views"
)

var (
	wsConn      *websocket.Conn
	currentView internal.AppView = internal.ViewResolver
)

func main() {
	go func() {

		if err := connectToWebSocket("ws://localhost:8080/ws"); err != nil {
			log.Println("WebSocket connection error:", err)
		}

		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func connectToWebSocket(wsURL string) error {
	u, err := url.Parse(wsURL)
	if err != nil {
		return err
	}

	log.Println("Connecting to", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	wsConn = c
	log.Println("WebSocket connected.")
	return nil
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
				views.LayoutMain(gtx, th, &currentView)
			case internal.ViewLogin:
				views.LayoutLogin(gtx, th, &currentView)
			case internal.ViewRegister:
				views.LayoutRegister(gtx, th, &currentView)
			case internal.ViewResolver:
				views.LayoutResolver(gtx, th, &currentView)
			case internal.ViewLoading:
				views.LayoutLoading(gtx, th, &currentView)
			case internal.ViewError:
				views.LayoutError(gtx, th, &currentView)
			case internal.ViewLogged:
				views.LayoutLogged(gtx, th, &currentView)
			}

			e.Frame(gtx.Ops)
		}
	}
}
