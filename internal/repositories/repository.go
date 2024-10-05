package repositories

// UserRepository - interface for user data managment
type UserRepository interface {
	StoreRefreshToken(userID, ip, email string, hashedToken []byte) error
	VerifyRefreshToken(userID, refreshToken string) (bool, error)
	VerifyIP(userID, ip string) (bool, string, error)
}
