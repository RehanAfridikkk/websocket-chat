package handlers

import (
	"fmt"

	structure "websocket-chat/struct"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan structure.MessageWithSender)

func HandleMessages() {
	for {
		msgWithSender := <-broadcast
		sender := msgWithSender.Sender
		msg := msgWithSender.Message

		if msg.To != "" {
			targetClient, found := findClientByUsername(msg.To, sender)
			if found && targetClient != nil {
				err := targetClient.WriteJSON(structure.Message{
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
				err := sender.WriteJSON(structure.Message{
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

				err := client.WriteJSON(structure.Message{
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
