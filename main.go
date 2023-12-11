package main

import (
	"fmt"
	"net/http"
	"websocket-chat/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/ws", handlers.HandleConnections)

	go handlers.HandleMessages()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
