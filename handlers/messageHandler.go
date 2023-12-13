package handlers

import (
	"fmt"

	structure "websocket-chat/struct"
	"websocket-chat/utils"

	"github.com/gorilla/websocket"
)

func HandleMessages() {
	for {
		msgWithSender := <-broadcast
		sender := msgWithSender.Sender
		msg := msgWithSender.Message

		if msg.To != "" {
			targetClient, found := FindClientByUsername(msg.To, sender)
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

				roomId := msg.Room
				db, _ := utils.OpenDB()

				var roomClients []structure.Client
				db.Where("room_id = ?", roomId).Find(&roomClients)

				for range roomClients {
					// Assuming structure.Client has no Sender field
					targetClient := client

					err := targetClient.WriteJSON(structure.Message{
						Username: msg.Username,
						Message:  msg.Message,
					})
					if err != nil {
						fmt.Println(err)
						targetClient.Close()
						// Optionally, remove the disconnected client from the roomClients slice.
						// This depends on your requirements.
					}
				}
			}
		}
	}
}

func FindClientByUsername(username string, sender *websocket.Conn) (*websocket.Conn, bool) {
	for client, u := range clients {
		if u == username && client != sender {
			return client, true
		}
	}
	return nil, false
}
