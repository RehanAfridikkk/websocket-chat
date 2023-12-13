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
			for client, clientUsername := range clients {
				if client == sender {
					continue
				}

				roomname := msg.RoomName
				db, _ := utils.OpenDB()

				var roomClients []structure.Client
				db.Where("room_name = ?", roomname).Find(&roomClients)

				for _, roomClient := range roomClients {

					if roomClient.Username == clientUsername {
						client.WriteJSON(structure.Message{
							Username: msg.Username,
							Message:  msg.Message,
						})
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
