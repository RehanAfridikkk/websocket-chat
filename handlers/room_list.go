package handlers

import (
	"log"
	structure "websocket-chat/struct"
	"websocket-chat/utils"

	"github.com/labstack/echo"
)

func RoomList(c echo.Context) error {

	var roomsList []structure.ChatRoom
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Find(&roomsList).Error
	if err != nil {
		log.Println(err)
		return c.JSON(500, echo.Map{"error": "Failed to retrieve room list"})
	}

	var roomDetails []map[string]interface{}
	for _, room := range roomsList {
		roomDetails = append(roomDetails, map[string]interface{}{
			"id":   room.ID,
			"name": room.Name,
		})
	}

	return c.JSON(200, echo.Map{"rooms": roomDetails})
}
