package services

import (
	"fmt"
	"medods/internal/config"
	"medods/internal/entities"
	"medods/internal/repositories"
	"medods/internal/utils"

	"github.com/google/uuid"
)

type authService struct {
	userRepo   repositories.UserRepository
	jwtUtils   utils.JWTUtils
	emailUtils utils.EmailUtils
	email      string
}

type AuthService interface {
	GenerateTokens(userID, ipAddress, email string) (*entities.TokenPair, error)
	RefreshTokens(accessToken, refreshToken string) (*entities.TokenPair, error)
}

// NewAuthService - create a new copy of service
func NewAuthService(jwtSecret string, userRepo repositories.UserRepository, emailSMTP config.EmailConfig) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtUtils:   utils.NewJWTUtils(jwtSecret),
		emailUtils: utils.NewEmailUtils(emailSMTP.SMTPHost, emailSMTP.SMTPUsername, emailSMTP.SMTPPassword, emailSMTP.SMTPPort),
		email:      emailSMTP.SMTPUsername + "@mail.com",
	}
}

// GenerateTokens - create a pair of refresh and access token
func (s *authService) GenerateTokens(userID, ipAddress, email string) (*entities.TokenPair, error) {
	accessToken, err := s.jwtUtils.GenerateAccessToken(userID, ipAddress)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	err = s.userRepo.StoreRefreshToken(userID, refreshToken, ipAddress, email)
	if err != nil {
		return nil, err
	}

	return &entities.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshTokens - realise a refresh operation
func (s *authService) RefreshTokens(acccessToken, refreshToken string) (*entities.TokenPair, error) {
	userClaims, err := s.jwtUtils.ParseJWT(acccessToken)
	if err != nil {
		return nil, err
	}

	valid, err := s.userRepo.VerifyRefreshToken(userClaims.ID, refreshToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("unvalid token")
	}

	ok, email, err := s.userRepo.VerifyIP(userClaims.ID, userClaims.IP)
	if err != nil {
		return nil, err
	}

	if !ok {
		emailMessage := entities.EmailRequest{
			From:    s.email,
			To:      email,
			Subject: "Is it you?",
			Body:    fmt.Sprintf("Is it your IP: %s", userClaims.IP),
		}
		s.emailUtils.SendEmail(emailMessage)
	}

	return s.GenerateTokens(userClaims.ID, userClaims.IP, email)
}
