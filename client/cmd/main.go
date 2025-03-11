package main

import (
	"client/internal/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := ui.CreateMainWindow(a)
	w.ShowAndRun()
}
