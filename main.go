package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


type MessageWithSender struct {
	Sender  *websocket.Conn
	Message Message
}

type Message struct {
	Username string `json:"username"`
	To       string `json:"to"`      
	Message  string `json:"message"`
}


var clients = make(map[*websocket.Conn]string) 
var broadcast = make(chan MessageWithSender)




func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Chat Room!")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	username := r.Header.Get("username")

	clients[conn] = username

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}

		broadcast <- MessageWithSender{Sender: conn, Message: msg}
	}
}



func handleMessages() {
	for {
		msgWithSender := <-broadcast
		sender := msgWithSender.Sender
		msg := msgWithSender.Message

		if msg.To != "" {
			targetClient, found := findClientByUsername(msg.To, sender)
			if found && targetClient != nil { 
				err := targetClient.WriteJSON(Message{
					Username: msg.Username,
					To:       msg.To,
					Message:  msg.Message,
				})
				if err != nil {
					fmt.Println(err)
					targetClient.Close()
					delete(clients, targetClient)
				}
			} else {
				err := sender.WriteJSON(Message{
					Username: "System",
					Message:  "User does not exist.",
				})
				if err != nil {
					fmt.Println(err)
					sender.Close()
					delete(clients, sender)
				}
			}
		} else {
			for client, _ := range clients {
				if client == sender {
					continue
				}

				err := client.WriteJSON(Message{
					Username: msg.Username,
					Message:  msg.Message,
				})
				if err != nil {
					fmt.Println(err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}



func findClientByUsername(username string, sender *websocket.Conn) (*websocket.Conn, bool) {
	for client, u := range clients {
		if u == username && client != sender {
			return client, true
		}
	}
	return nil, false
}


