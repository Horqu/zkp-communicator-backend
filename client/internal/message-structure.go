package internal

type Message struct {
	Command MessageCommand `json:"command"`
	Data    string         `json:"data"`
}

type Response struct {
	Command ResponseCommand `json:"command"`
	Data    string          `json:"data"`
}
