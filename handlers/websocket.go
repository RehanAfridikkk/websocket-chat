package handlers

import (
	"log"
	"sync"
	structure "websocket-chat/struct"
	"websocket-chat/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan structure.MessageWithSender)
	mu        sync.Mutex
)

func HandleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	senderUsername, err := utils.ExtractUsernameFromToken(c)
	if err != nil {
		log.Println(err)
		return err
	}

	mu.Lock()
	clients[conn] = senderUsername
	mu.Unlock()

	for {
		var msg structure.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return err
		}

		msg.Username = senderUsername

		broadcast <- structure.MessageWithSender{Sender: conn, Message: msg}
	}
}
