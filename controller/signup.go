package controller

import (
	"fmt"
	"log"
	"net/http"
	"websocket-chat/models"
	"websocket-chat/utils"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type SignUpInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signup(c echo.Context) error {

	input := new(SignUpInput)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := validateUser(*input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		fmt.Println("error", err)
	}
	user := models.User{
		Username: input.Username,
		Password: hashedPassword,
		Role:     "User",
	}
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Create(&user).Error
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusConflict, "User with the provided username already exists")

	}

	return c.JSON(http.StatusCreated, user)
}
func validateUser(user SignUpInput) error {

	if user.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username cannot be empty")
	}
	if len(user.Password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
