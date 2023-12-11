package structure

import "github.com/gorilla/websocket"

type MessageWithSender struct {
	Sender  *websocket.Conn
	Message Message
}

type Message struct {
	Username string `json:"username"`
	To       string `json:"to"`
	Message  string `json:"message"`
}
