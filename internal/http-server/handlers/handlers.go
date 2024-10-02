package handlers

import (
	"errors"
	"log/slog"
	"medods/internal/http-server/models"
	"medods/internal/infrustructure/logger/sl"
	"medods/internal/services"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
)

// authHandlers - structure, which realize AuthHandlers interface
type authHandlers struct {
	logger   *slog.Logger
	services services.AuthService
}

// AuthHandlers - interface, which realize our API
type AuthHandlers interface {
	AccessHandler(ctx *gin.Context)
	RefreshHanadler(ctx *gin.Context)
}

// NewAuthHandlers - create a new copy of AuthHandlers interface
func NewAuthHandlers(logger *slog.Logger, services services.AuthService) AuthHandlers {
	return &authHandlers{
		logger:   logger,
		services: services,
	}
}

// AccessHandler - is a handler which generate new pair of access and refresh token and put them into cookie
func (h *authHandlers) AccessHandler(ctx *gin.Context) {
	var req models.AccessRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", sl.Err(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		h.logger.Error("Invalid email", "email", req.Email)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	clientIP := ctx.ClientIP()

	pairToken, err := h.services.GenerateTokens(req.GUID, clientIP, req.Email)
	if err != nil {
		h.logger.Error("failed generate a token pair", sl.Err(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	response := models.Response{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
	}

	ctx.JSON(http.StatusOK, response)
}

// RefreshHandler - is a hadnler, which check a refresh token and if all is ok return token pair
func (h *authHandlers) RefreshHanadler(ctx *gin.Context) {
	var req models.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", sl.Err(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bab request",
		})
		return
	}

	pairToken, err := h.services.RefreshTokens(req.AccessToken, req.RefreshToken)
	if err != nil {
		if errors.Is(err, errors.New("invalid token")) {
			h.logger.Warn("Expired refresh token")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "expired refresh token",
			})
			return
		}

		h.logger.Error("failed to refresh tokens", sl.Err(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}

	responce := models.Response{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
	}

	ctx.JSON(http.StatusOK, responce)
}
