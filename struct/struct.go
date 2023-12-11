package structure 


type MessageWithSender struct {
	Sender  *websocket.Conn
	Message Message
}
