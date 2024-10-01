package entities

import "github.com/dgrijalva/jwt-go"

// User struct, which represent our user
type User struct {
	ID          string
	HashedToken string
	IP          string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserClaims struct {
	ID string `json:"user_id"`
	IP string `json:"ip"`
	jwt.StandardClaims
}
