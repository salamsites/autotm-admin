package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TokenEncode(userID int64, role string, jwtKey string) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
