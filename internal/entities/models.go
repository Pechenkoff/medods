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
	jwt.StandardClaims
}

type EmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
