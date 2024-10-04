package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"medods/internal/entities"
	"medods/internal/http-server/handlers"
	"medods/internal/http-server/handlers/mocks"
	"medods/internal/http-server/models"
	"medods/internal/infrustructure/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		body             models.AccessRequest
		mockReturn       *entities.TokenPair
		mockError        error
		expectedCode     int
		expectedResponse interface{}
	}{
		{
			name:             "Success",
			body:             models.AccessRequest{GUID: "test-guid", Email: "test@example.com"},
			mockReturn:       &entities.TokenPair{AccessToken: "access-token", RefreshToken: "refresh-token"},
			mockError:        nil,
			expectedCode:     http.StatusOK,
			expectedResponse: map[string]interface{}{"access_token": "access-token", "refresh_token": "refresh-token"},
		},
		{
			name:             "bad email format",
			body:             models.AccessRequest{GUID: "test-guid", Email: "invalid-email"},
			mockReturn:       nil,
			mockError:        nil,
			expectedCode:     http.StatusBadRequest,
			expectedResponse: map[string]interface{}{"error": "invalid email"},
		},
		{
			name:             "Server error",
			body:             models.AccessRequest{GUID: "test-guid", Email: "test@example.com"},
			mockReturn:       nil,
			mockError:        errors.New("service error"),
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{"error": "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewAuthService(t)
			if tt.name != "bad email format" {
				mockService.On("GenerateTokens", tt.body.GUID, mock.Anything, tt.body.Email).Return(tt.mockReturn, tt.mockError)
			}

			handler := handlers.NewAuthHandlers(slogdiscard.NewDiscardLogger(), mockService)

			router := gin.New()
			router.POST("/access", handler.AccessHandler)

			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/access", bytes.NewBuffer(jsonBody))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed bind response: %v", err)
			}
			assert.Equal(t, tt.expectedResponse, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		body             models.RefreshRequest
		mockReturn       *entities.TokenPair
		mockError        error
		expectedCode     int
		expectedResponse interface{}
	}{
		{
			name:             "Success",
			body:             models.RefreshRequest{AccessToken: "access_token", RefreshToken: "refresh_token"},
			mockReturn:       &entities.TokenPair{AccessToken: "access_token", RefreshToken: "refresh_token"},
			mockError:        nil,
			expectedCode:     http.StatusOK,
			expectedResponse: map[string]interface{}{"access_token": "access_token", "refresh_token": "refresh_token"},
		},
		{
			name:             "bad token",
			body:             models.RefreshRequest{AccessToken: "access_token", RefreshToken: "refresh_token"},
			mockReturn:       nil,
			mockError:        errors.New("invalid token"),
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{"error": "expired refresh token"},
		},
		{
			name:             "server error",
			body:             models.RefreshRequest{AccessToken: "access_token", RefreshToken: "refresh_token"},
			mockReturn:       nil,
			mockError:        errors.New("server error"),
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{"error": "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewAuthService(t)

			mockService.On("RefreshTokens", tt.body.AccessToken, tt.body.RefreshToken, mock.Anything).Return(tt.mockReturn, tt.mockError)

			handler := handlers.NewAuthHandlers(slogdiscard.NewDiscardLogger(), mockService)

			router := gin.New()
			router.POST("/refresh", handler.RefreshHanadler)

			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(jsonBody))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed bind response: %v", err)
			}
			assert.Equal(t, tt.expectedResponse, response)

			mockService.AssertExpectations(t)
		})
	}
}
