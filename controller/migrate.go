package controller

import (
	"fmt"
	"websocket-chat/models"
	"websocket-chat/utils"

	"github.com/labstack/echo"
)

func Migrate(c echo.Context) error {
	db, err := utils.OpenDB()
	if err != nil {
		fmt.Println("Unable to connect to DB")
	}

	db.AutoMigrate(&models.User{})
	return nil
}
