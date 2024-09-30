package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateAccessToken(userID, ip, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"ip":      ip,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secret))
}
