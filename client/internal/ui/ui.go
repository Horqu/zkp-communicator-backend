package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateMainWindow sets up the main window layout for the chat client application.
func CreateMainWindow(app fyne.App) fyne.Window {
	mainWindow := app.NewWindow("Chat Client")

	// Create input field for message
	messageEntry := widget.NewEntry()
	messageEntry.SetPlaceHolder("Type your message...")

	// Create send button
	sendButton := widget.NewButton("Send", func() {
		// Logic to send the message will be implemented here
	})

	// Create a layout with the message entry and send button
	content := container.NewVBox(
		messageEntry,
		sendButton,
	)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(400, 300))
	mainWindow.Show()

	return mainWindow
}