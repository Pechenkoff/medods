package services_test

import (
	"errors"
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

}
