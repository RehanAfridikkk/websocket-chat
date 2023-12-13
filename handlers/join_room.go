package handlers

import (
	"log"
	structure "websocket-chat/struct"
	"websocket-chat/utils"

	"github.com/labstack/echo"
)

func JoinRoom(c echo.Context) error {
	// Extract room name and password from the request
	roomName := c.FormValue("name")
	password := c.FormValue("password")

	// Retrieve the room from the database
	var room structure.ChatRoom
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Where("name = ? AND password = ?", roomName, password).First(&room).Error
	if err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid room name or password"})
	}

	// Extract sender's username from JWT in Authorization header
	senderUsername, err := utils.ExtractUsernameFromToken(c)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create a new client record for the user joining the room
	client := structure.Client{
		Username: senderUsername,
		RoomID:   room.ID,
		RoomName: room.Name,
	}

	// Save the client record in the database
	err = db.Create(&client).Error
	if err != nil {
		log.Println(err)
		return c.JSON(500, echo.Map{"error": "Failed to join room"})
	}

	return c.JSON(200, echo.Map{"message": "Joined room successfully"})
}
