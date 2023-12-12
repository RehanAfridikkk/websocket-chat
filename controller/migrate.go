package controller

import (
	"fmt"
	"log"
	"websocket-chat/models"
	structure "websocket-chat/struct"
	"websocket-chat/utils"

	"github.com/labstack/echo"
)

func Migrate(c echo.Context) error {
	db, err := utils.OpenDB()
	if err != nil {
		fmt.Println("Unable to connect to DB")
		return err
	}

	// Create ChatRoom and Client tables
	err = db.AutoMigrate(&structure.ChatRoom{}, &structure.Client{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Set up foreign key relationship
	err = db.AutoMigrate(&structure.ChatRoom{}, &structure.Client{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	// AutoMigrate other models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
