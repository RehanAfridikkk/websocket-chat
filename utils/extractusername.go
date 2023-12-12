package utils

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
)

func ExtractUsernameFromToken(c echo.Context) (string, error) {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("authorization token not provided")
	}

	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", fmt.Errorf("invalid token format")
	}

	token, err := jwt.Parse(tokenParts[1], func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	username, ok := claims["name"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in token claims")
	}

	return username, nil
}
