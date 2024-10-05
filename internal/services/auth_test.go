package services_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"medods/internal/entities"
	"medods/internal/services"
	"medods/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name           string
		mockJWTReturn  string
		mockJWTError   error
		mockRepoError  error
		expectedReturn *entities.TokenPair
		expectedError  error
	}{
		{
			name:           "Success",
			mockJWTReturn:  "access_token",
			mockJWTError:   nil,
			mockRepoError:  nil,
			expectedReturn: &entities.TokenPair{AccessToken: "access_token", RefreshToken: string(mock.AnythingOfType("string"))},
			expectedError:  nil,
		},
		{
			name:           "bad jwt",
			mockJWTReturn:  "",
			mockJWTError:   errors.New("couldn't create roken"),
			mockRepoError:  nil,
			expectedReturn: nil,
			expectedError:  errors.New("couldn't create roken"),
		},
		{
			name:           "failed in repo",
			mockJWTReturn:  "access_token",
			mockJWTError:   nil,
			mockRepoError:  errors.New("rewrite query failed"),
			expectedReturn: nil,
			expectedError:  errors.New("rewrite query failed"),
		},
	}

	for _, tt := range tests {
		mockJWTUtils := mocks.NewJWTUtils(t)
		mockProducer := mocks.NewProducer(t)
		mockRepo := mocks.NewUserRepository(t)
		if tt.name != "bad jwt" {
			mockRepo.On("StoreRefreshToken", "user_id", mock.AnythingOfType("string"), "ip", "email").Return(tt.mockRepoError)
		}
		mockJWTUtils.On("GenerateAccessToken", "user_id").Return(tt.mockJWTReturn, tt.mockJWTError)
		service := services.NewAuthService(mockJWTUtils, mockRepo, mockProducer)
		tokenPair, err := service.GenerateTokens("user_id", "ip", "email")

		if tt.expectedReturn != nil {
			assert.Equal(t, tt.expectedReturn.AccessToken, tokenPair.AccessToken)
		}

		assert.Equal(t, tt.expectedError, err)
	}
}

func TestRefreshTokens(t *testing.T) {
	tests := []struct {
		name                  string
		mockJWTUtilsReturn    *entities.UserClaims
		mockJWTUtilsError     error
		mockVerifyTokenReturn bool
		mockVerifyTokenError  error
		mockVerifyIPReturn    bool
		mockVerifyIPString    string
		mockVerifyIPError     error
		mockSendMessageError  error
		expectedReturn        *entities.TokenPair
		expectedError         error
	}{
		{
			name:                  "Success, but not pass verify IP",
			mockJWTUtilsReturn:    &entities.UserClaims{ID: "user_id"},
			mockJWTUtilsError:     nil,
			mockVerifyTokenReturn: true,
			mockVerifyTokenError:  nil,
			mockVerifyIPReturn:    false,
			mockVerifyIPString:    "email",
			mockVerifyIPError:     nil,
			mockSendMessageError:  nil,
			expectedReturn:        &entities.TokenPair{AccessToken: "access_token", RefreshToken: string(mock.AnythingOfType("string"))},
			expectedError:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWTUtils := mocks.NewJWTUtils(t)
			mockProducer := mocks.NewProducer(t)
			mockRepo := mocks.NewUserRepository(t)
			mockRepo.On("VerifyRefreshToken", "user_id", mock.Anything).Return(tt.mockVerifyTokenReturn, tt.mockVerifyTokenError)
			mockRepo.On("VerifyIP", "user_id", "ip").Return(tt.mockVerifyIPReturn, tt.mockVerifyIPString, tt.mockVerifyIPError)

			message := services.Message{
				Email:   "email",
				Subject: "Is it you?",
				Message: fmt.Sprintf("We are see that some one trying to login with your credentials from this ip: %v, is it you?", "ip"),
			}

			msgByte, _ := json.Marshal(message)

			mockProducer.On("SendMessage", "change ip email", msgByte).Return(tt.mockSendMessageError)

			mockJWTUtils.On("ParseJWT", "access_token").Return(tt.mockJWTUtilsReturn, tt.mockJWTUtilsError)

			mockJWTUtils.On("GenerateAccessToken", "user_id").Return(tt.expectedReturn.AccessToken, nil)

			mockRepo.On("StoreRefreshToken", "user_id", mock.AnythingOfType("string"), "ip", "email").Return(nil)

			service := services.NewAuthService(mockJWTUtils, mockRepo, mockProducer)
			tokenpair, err := service.RefreshTokens("access_token", "refresh_token", "ip")

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedReturn.AccessToken, tokenpair.AccessToken)
		})
	}
}
