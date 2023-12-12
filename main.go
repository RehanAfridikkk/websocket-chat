package main

import (
	"log"
	"websocket-chat/controller"
	"websocket-chat/handlers"
	"websocket-chat/utils"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	// http.HandleFunc("/", handlers.HomePage)
	// http.HandleFunc("/ws", handlers.HandleConnections)

	go handlers.HandleMessages()

	e.POST("/signup", controller.Signup)
	e.POST("/login", controller.Login)
	e.POST("/migrate", controller.Migrate)
	e.GET("/ws", handlers.HandleWebSocket)
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	e.POST("/create-room", echo.HandlerFunc(func(c echo.Context) error {
		return handlers.CreateRoom(c, db)
	}))
	e.POST("/join-room", handlers.JoinRoom)
	e.GET("/room-list", handlers.RoomList)
	// fmt.Println("Server started on :8080")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	panic("Error starting server: " + err.Error())
	// }

	// Defer closing the database connection
	e.Logger.Fatal(e.Start(":1304"))
	// Pass the gorm.DB instance to the controller
	// controller.SetDB(db)

}
