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

	return c.JSON(200, echo.Map{"message": "Joined room successfully"})
}
