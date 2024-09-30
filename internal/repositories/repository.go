package repositories

// UserRepository - interface for user data managment
type UserRepository interface {
	StoreRefreshToken(userID, refreshToken, ip string) error
	VerifyRefreshToken(userID, refreshToken string) (bool, error)
	VerifyIP(userID, ip string) (bool, error)
}
