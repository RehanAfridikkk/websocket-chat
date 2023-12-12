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
	// Extract room name and password from the request
	roomName := c.FormValue("name")
	password := c.FormValue("password")

	senderUsername, err := utils.ExtractUsernameFromToken(c)
	if err != nil {
		log.Println(err)
		return err
	}
	// Check if the room already exists
	var existingRoom structure.ChatRoom
	err = db.Where("name = ?", roomName).First(&existingRoom).Error
	if err == nil {
		return c.JSON(400, echo.Map{"error": "Room already exists"})
	}

	// Create a new chat room
	newRoom := structure.ChatRoom{
		Name:     roomName,
		Password: password,
	}

	// Store the room in the database
	err = db.Create(&newRoom).Error
	if err != nil {
		log.Println(err)
		return c.JSON(500, echo.Map{"error": "Failed to create room"})
	}

	// Extract sender's username from JWT in Authorization header

	// Retrieve the created room from the database
	err = db.Where("name = ?", roomName).First(&newRoom).Error
	if err != nil {
		log.Println(err)
		return c.JSON(400, echo.Map{"error": "Invalid room name"})
	}

	// Create a client record for the user who created the room
	mu.Lock()
	client := structure.Client{
		Username: senderUsername, // You can set this to nil or a default value based on your use case
		RoomID:   newRoom.ID,
	}
	db.Create(&client)
	mu.Unlock()

	return c.JSON(200, echo.Map{"message": "Room created successfully"})
}
