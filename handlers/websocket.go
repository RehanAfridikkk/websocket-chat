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

	roomName := c.Param("room")

	// Retrieve the chat room from the database
	var room structure.ChatRoom
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Where("name = ?", roomName).First(&room).Error
	if err != nil {
		log.Println(err)
		return c.JSON(400, echo.Map{"error": "Invalid room name"})
	}

	// Extract username from JWT in Authorization header
	senderUsername, err := utils.ExtractUsernameFromToken(c)
	if err != nil {
		log.Println(err)
		return err
	}

	mu.Lock()
	client := structure.Client{
		Username: senderUsername,

		RoomID: room.ID,
	}
	room.Clients = append(room.Clients, client)
	clientsMap[conn] = &room
	mu.Unlock()

	for {
		var msg structure.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			mu.Lock()
			// Remove the disconnected client from the slice
			updatedClients := make([]structure.Client, 0, len(room.Clients)-1)
			for _, c := range room.Clients {
				if c.Username != senderUsername {
					updatedClients = append(updatedClients, c)
				}
			}
			room.Clients = updatedClients
			mu.Unlock()
			return err
		}

		msg.Username = senderUsername
		msg.Room = roomName

		broadcast <- structure.MessageWithSender{Sender: conn, Message: msg}
	}
}
