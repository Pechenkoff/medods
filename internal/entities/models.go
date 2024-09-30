package entities

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
