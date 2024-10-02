package services

import (
	"encoding/json"
	"fmt"
	"medods/internal/entities"
	kafka "medods/internal/infrustructure/kafka/producer"
	"medods/internal/repositories"
	"medods/internal/utils"

	"github.com/google/uuid"
)

// Message - structure of message to kafka
type Message struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// authService is a structure, which help to realize a AuthService interface
type authService struct {
	userRepo repositories.UserRepository
	jwtUtils utils.JWTUtils
	producer kafka.Producer
}

// AuthService is a interface, which presents auth service
type AuthService interface {
	GenerateTokens(userID, ipAddress, email string) (*entities.TokenPair, error)
	RefreshTokens(accessToken, refreshToken string) (*entities.TokenPair, error)
}

// NewAuthService - create a new copy of service
func NewAuthService(jwtSecret string, userRepo repositories.UserRepository, producer kafka.Producer) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtUtils: utils.NewJWTUtils(jwtSecret),
		producer: producer,
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
		return nil, fmt.Errorf("invalid token")
	}

	ok, email, err := s.userRepo.VerifyIP(userClaims.ID, userClaims.IP)
	if err != nil {
		return nil, err
	}

	if !ok {
		message := Message{
			Email:   email,
			Subject: "Is it you?",
			Message: fmt.Sprintf("We are see that some one trying to login with your credentials from this ip: %v, is it you?", userClaims.IP),
		}

		msgByte, err := json.Marshal(message)
		if err != nil {
			return nil, err
		}

		err = s.producer.SendMessage("change ip email", msgByte)
		if err != nil {
			return nil, err
		}
	}

	return s.GenerateTokens(userClaims.ID, userClaims.IP, email)
}
