package handlers

import (
	"log"
	structure "websocket-chat/struct"
	"websocket-chat/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

var (
	rooms = make(map[string]*structure.ChatRoom)

	clientsMap = make(map[*websocket.Conn]*structure.ChatRoom)
)

func CreateRoom(c echo.Context, db *gorm.DB) error {
	roomName := c.FormValue("name")
	password := c.FormValue("password")

	senderUsername, err := utils.ExtractUsernameFromToken(c)
	if err != nil {
		log.Println(err)
		return err
	}
	var existingRoom structure.ChatRoom
	err = db.Where("name = ?", roomName).First(&existingRoom).Error
	if err == nil {
		return c.JSON(400, echo.Map{"error": "Room already exists"})
	}

	newRoom := structure.ChatRoom{
		Name:     roomName,
		Password: password,
	}

	err = db.Create(&newRoom).Error
	if err != nil {
		log.Println(err)
		return c.JSON(500, echo.Map{"error": "Failed to create room"})
	}

	err = db.Where("name = ?", roomName).First(&newRoom).Error
	if err != nil {
		log.Println(err)
		return c.JSON(400, echo.Map{"error": "Invalid room name"})
	}

	mu.Lock()
	client := structure.Client{
		Username: senderUsername,
		RoomID:   newRoom.ID,
		RoomName: newRoom.Name,
	}
	db.Create(&client)
	mu.Unlock()

	return c.JSON(200, echo.Map{"message": "Room created successfully"})
}
