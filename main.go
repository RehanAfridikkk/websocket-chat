package main

import (
	"websocket-chat/controller"
	"websocket-chat/handlers"

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
