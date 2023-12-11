package structure

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

type MessageWithSender struct {
	Sender  *websocket.Conn
	Message Message
}

type Message struct {
	Username string `json:"username"`
	To       string `json:"to"`
	Message  string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

// type User struct {
// 	gorm.Model

// 	Username string `gorm:"unique;not null"`
// 	Password string `gorm:"not null"`
// 	Role     string `gorm:"not null"`
// }
