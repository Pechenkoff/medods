package services

import (
	"fmt"
	"medods/internal/entities"
	"medods/internal/repositories"
	"medods/internal/utils"

	"github.com/google/uuid"
)

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

type AuthService interface {
	GenerateTokens(userID, ipAddress string) (*entities.TokenPair, error)
	RefreshTokens(userID, oldRefreshToken, newIP string) (*entities.TokenPair, error)
}

// NewAuthService - create a new copy of service
func NewAuthService(jwtSecret string, userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// GenerateTokens - create a pair of refresh and access token
func (s *authService) GenerateTokens(userID, ipAddress string) (*entities.TokenPair, error) {
	accessToken, err := utils.GenerateAccessToken(userID, ipAddress, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	err = s.userRepo.StoreRefreshToken(userID, refreshToken, ipAddress)
	if err != nil {
		return nil, err
	}

	return &entities.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshTokens - realise a refresh operation
func (s *authService) RefreshTokens(userID, oldRefreshToken, ipAddress string) (*entities.TokenPair, error) {
	valid, err := s.userRepo.VerifyRefreshToken(userID, oldRefreshToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("uncalid token")
	}

	return s.GenerateTokens(userID, ipAddress)
}
