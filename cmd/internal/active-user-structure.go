package internal

import (
	"time"

	"github.com/gorilla/websocket"
)

type ActiveUser struct {
	WsConnection   *websocket.Conn
	Expiry         time.Time
	PublicKey      string
	SelectedFriend string
}
