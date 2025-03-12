package views

import (
	"fmt"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"client/internal"
)

// Zmienne do logiki ekranu Logged
var (
	// Czas sesji w sekundach (2 minuty = 120s)
	sessionTimeLeft = 120
	lastFrameLogged time.Time

	// Górny pasek
	refreshButton = new(widget.Clickable)
	logoutButton  = new(widget.Clickable)

	// Lista znajomych i zaznaczenie aktualnego
	friendList     = []string{"Alice", "Bob", "Charlie"}
	selectedFriend = -1

	// Dolny pasek: dodawanie znajomego i wysyłanie wiadomości
	newFriendEditor   = new(widget.Editor)
	addFriendButton   = new(widget.Clickable)
	messageEditor     = new(widget.Editor)
	sendMessageButton = new(widget.Clickable)

	friendButtons []*widget.Clickable
)

// LayoutLogged - główny ekran (chat) po zalogowaniu
func LayoutLogged(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
	now := time.Now()
	if !lastFrameLogged.IsZero() {
		dt := now.Sub(lastFrameLogged).Seconds()
		if dt >= 1 {
			// Zliczaj pełne sekundy
			secondsPassed := int(dt)
			sessionTimeLeft -= secondsPassed
			if sessionTimeLeft < 0 {
				sessionTimeLeft = 0
			}
		}
	}
	lastFrameLogged = now

	// Układ pionowy (top bar, środek, bottom bar)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Górny pasek
			return layoutTopBar(gtx, th, currentView)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Środek aplikacji
			return layoutMiddle(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Dolny pasek
			return layoutBottomBar(gtx, th)
		}),
	)
}

// layoutTopBar tworzy pasek z "session time left", przyciskiem Refresh i Logout
func layoutTopBar(gtx layout.Context, th *material.Theme, currentView *internal.AppView) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Sesja - format "mm:ss"
			minutes := sessionTimeLeft / 60
			seconds := sessionTimeLeft % 60
			sessionLabel := material.Label(th, unit.Sp(16), fmt.Sprintf("Session time left: %d:%02d", minutes, seconds))
			sessionLabel.Alignment = text.Start
			sessionLabel.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
			return sessionLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk Refresh
			btn := material.Button(th, refreshButton, "Refresh")
			if refreshButton.Clicked(gtx) {
				// Tu można odświeżyć listę znajomych, wiadomości, etc.
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk Logout
			btn := material.Button(th, logoutButton, "Logout")
			if logoutButton.Clicked(gtx) {
				// Powrót do ekranu logowania
				*currentView = internal.ViewMain
				sessionTimeLeft = 120 // resetujemy ewentualnie czas sesji
			}
			return btn.Layout(gtx)
		}),
	)
}

// layoutMiddle - środkowa część aplikacji: lista znajomych (lewa) i chat z wybranym znajomym (prawa)
func layoutMiddle(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Lista znajomych, pionowo
			return layoutFriendsList(gtx, th)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Chat z aktualnie zaznaczonym znajomym
			return layoutChat(gtx, th)
		}),
	)
}

func layoutFriendsList(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Jeśli friendButtons jest puste lub długość się nie zgadza, inicjalizujemy ponownie
	if len(friendButtons) != len(friendList) {
		friendButtons = make([]*widget.Clickable, len(friendList))
		for i := range friendButtons {
			friendButtons[i] = new(widget.Clickable)
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		friendItems(gtx, th)...,
	)
}

func friendItems(gtx layout.Context, th *material.Theme) []layout.FlexChild {
	var children []layout.FlexChild
	for i, friend := range friendList {
		index := i
		button := friendButtons[i] // Pobieramy już istniejący widget.Clickable

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			item := material.Button(th, button, friend)
			if button.Clicked(gtx) {
				fmt.Printf("Selected friend: %s\n", friend)
				selectedFriend = index
			}

			// Podświetlenie wybranego znajomego
			if index == selectedFriend {
				item.Background = color.NRGBA{R: 0x00, G: 0xAA, B: 0xFF, A: 0xFF}
				item.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
			} else {
				item.Background = color.NRGBA{R: 0xBB, G: 0xBB, B: 0xBB, A: 0xFF}
				item.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
			}

			item.Inset = layout.UniformInset(unit.Dp(4))
			return item.Layout(gtx)
		}))
	}
	return children
}

// layoutChat wyświetla chat z aktualnie wybraną osobą
func layoutChat(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if selectedFriend < 0 || selectedFriend >= len(friendList) {
		// Jeśli nic nie wybrano, pokaż placeholder
		lbl := material.Label(th, unit.Sp(16), "Select a friend to chat")
		lbl.Alignment = text.Middle
		return lbl.Layout(gtx)
	}
	friendName := friendList[selectedFriend]
	lbl := material.Label(th, unit.Sp(16), "Chat with "+friendName)
	lbl.Alignment = text.Middle
	return lbl.Layout(gtx)
}

// layoutBottomBar - pole do wpisania nowego znajomego i wiadomości
func layoutBottomBar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Pole do wpisania nowego znajomego
			edit := material.Editor(th, newFriendEditor, "Add new friend")
			// edit.SingleLine = true
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk "Add friend"
			btn := material.Button(th, addFriendButton, "Add")
			if addFriendButton.Clicked(gtx) {
				// Dodaj znajomego do listy
				name := newFriendEditor.Text()
				if name != "" {
					friendList = append(friendList, name)
					newFriendEditor.SetText("")
				}
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Pole do wpisania wiadomości
			edit := material.Editor(th, messageEditor, "Type message")
			// edit.SingleLine = true
			return edit.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Przycisk "Send" do aktualnie wybranego znajomego
			btn := material.Button(th, sendMessageButton, "Send")
			if sendMessageButton.Clicked(gtx) {
				if selectedFriend >= 0 && selectedFriend < len(friendList) {
					msg := messageEditor.Text()
					// Tu możesz dodać logikę wysyłania wiadomości do friendList[selectedFriend]
					fmt.Printf("Send to %s: %s\n", friendList[selectedFriend], msg)
					messageEditor.SetText("")
				}
			}
			return btn.Layout(gtx)
		}),
	)
}
