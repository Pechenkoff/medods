package postgres

import (
	"context"
	"medods/internal/repositories"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *pgx.Conn
}

// NewUserRepository - create a new copy, which realese UserRepository interface
func NewUserRepository(db *pgx.Conn) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

// StoreRefreshToken - create a new user token and save it in PostgreSQL
func (r *userRepository) StoreRefreshToken(userID, ip, email string, hashedToken []byte) error {
	query := `INSERT INTO users (id, email, hashed_token, ip) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(context.Background(), query, userID, email, hashedToken, ip)
	if err != nil {
		return nil
	}

	return nil
}

// VerifyRefreshToken - check a refresh token
func (r *userRepository) VerifyRefreshToken(userID, refreshToken string) (bool, error) {
	var hashedToken string
	query := `SELECT hashed_token FROM users WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&hashedToken)
	if err != nil {
		return false, err
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refreshToken)) == nil, nil
}

// VerifyIP - check a IP, if not true return email, else return empty string
func (r *userRepository) VerifyIP(userID, ip string) (bool, string, error) {
	var ipOld, email string
	query := `SELECT ip, email FROM users WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&ipOld, &email)
	if err != nil {
		return false, "", err
	}

	return ip == ipOld, email, nil
}
