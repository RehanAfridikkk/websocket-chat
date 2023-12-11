// controller/user.go

package controller

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"websocket-chat/models"
	structure "websocket-chat/struct"
	"websocket-chat/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c echo.Context) error {
	loginRequest := new(structure.LoginRequest)
	if err := c.Bind(loginRequest); err != nil {
		log.Println("Error binding login request:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	username := loginRequest.Username
	password := loginRequest.Password

	log.Println("Login request:", loginRequest)

	db, _ := utils.OpenDB()

	user, err := GetUserByusername(db, loginRequest.Username)
	if err != nil {
		log.Println("Error checking credentials:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking credentials")
	}

	if user == nil {
		log.Println("User not found for ID:", loginRequest.Username)
		return echo.ErrUnauthorized
	}

	log.Println("Found user:", user)

	if !comparePasswords(password, user.Password) {
		log.Println("Password mismatch for user:", username)
		return echo.ErrUnauthorized
	}

	role := "admin"
	if user.Role != "admin" {
		role = "user"
	}

	refreshClaims := &structure.JwtCustomClaims{
		Name:  username,
		Admin: role == "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	rt, err := refreshToken.SignedString([]byte("refresh_secret"))
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating refresh token")
	}

	claims := &structure.JwtCustomClaims{
		Name:  username,
		Admin: role == "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 50)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println("Error generating token:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	log.Println("Login successful for user:", username)

	return c.JSON(http.StatusOK, echo.Map{
		"token":         t,
		"refresh_token": rt,
	})
}

func comparePasswords(providedPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		log.Println("Password comparison error:", err)
	}
	return err == nil
}
func GetUserByusername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.Model(models.User{}).Where("username = ?", username).Take(&user).Error
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println("Error querying user by ID:", err)
		return nil, err
	}

	return &user, nil
}
