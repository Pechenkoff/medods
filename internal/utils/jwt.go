package utils

import (
	"fmt"
	"medods/internal/entities"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtUtils struct {
	secret string
}

type JWTUtils interface {
	GenerateAccessToken(userID, ip string) (string, error)
	ParseJWT(tokenString string) (*entities.UserClaims, error)
}

func NewJWTUtils(secret string) JWTUtils {
	return &jwtUtils{
		secret: secret,
	}
}

func (u *jwtUtils) GenerateAccessToken(userID, ip string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"ip":      ip,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(u.secret))
}

func (u *jwtUtils) ParseJWT(tokenString string) (*entities.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entities.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}

		return []byte(u.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*entities.UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
