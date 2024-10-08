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
	GenerateAccessToken(userID string) (string, error)
	ParseJWT(tokenString string) (*entities.UserClaims, error)
}

// NewJWTUtils - create a new copy of JWTUtils
func NewJWTUtils(secret string) JWTUtils {
	return &jwtUtils{
		secret: secret,
	}
}

// GenerateAccessToken - generate a new JWT token which is depend of user_id and ip
func (u *jwtUtils) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(u.secret))
}

// ParseJWT - parse information from jwt
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
