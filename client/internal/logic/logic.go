package logic

import (
	"fmt"
	"net"
)

// ConnectToServer establishes a connection to the chat server.
func ConnectToServer(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	return conn, nil
}

// SendMessage sends a message to the connected chat server.
func SendMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}