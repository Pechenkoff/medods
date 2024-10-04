package handlers

import (
	"fmt"
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

// AccessHandler godoc
// @Summary Generate pair of access and refresh tokens by GUID and email
// @Description Retern a pair of access and refresh tokens by GUID and email
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.AccessRequest true "Body of request with GUID and email"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse400
// @Failure 500 {object} models.ErrorResponse500
// @Router /access [post]
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

// RefreshHandler godoc
// @Summary Get a new pair of access and refresh tokens
// @Description Generate a pair of access and refresh token, decifer access token get user GUID, check a refresh token and if all is ok return pair of access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RefreshRequest true "Body request with access and refresh token"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse400
// @Failure 500 {object} models.ErrorResponse500
// @Router /refresh [post]
func (h *authHandlers) RefreshHanadler(ctx *gin.Context) {
	var req models.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", sl.Err(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bab request",
		})
		return
	}

	userIP := ctx.ClientIP()

	pairToken, err := h.services.RefreshTokens(req.AccessToken, req.RefreshToken, userIP)
	if err != nil {
		if err.Error() == fmt.Errorf("invalid token").Error() {
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
		return
	}

	responce := models.Response{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
	}

	ctx.JSON(http.StatusOK, responce)
}
