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

// StroreRefreshToken - create a new user token and save it in PostgreSQL
func (r *userRepository) StoreRefreshToken(userID, ip, refreshToken string) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (id, hashed_token, ip) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(context.Background(), query, userID, hashedToken, ip)
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

// VerifyIP - check a IP
func (r *userRepository) VerifyIP(userID, ip string) (bool, error) {
	var ipOld string
	query := `SELECT ip FROM users WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&ipOld)
	if err != nil {
		return false, err
	}

	return ipOld == ip, nil
}
