package handlers

import (
	"fmt"
	"net/http"
	structure "websocket-chat/struct"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	username := r.Header.Get("username")

	clients[conn] = username

	for {
		var msg structure.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}

		broadcast <- structure.MessageWithSender{Sender: conn, Message: msg}
	}
}
