package handler

import (
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/storage/cache"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/token"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler interface {
	LoginHandler(c *gin.Context)
	IsAuthHandler(c *gin.Context)
	LogoutHandler(c *gin.Context)
	RefreshAccessTokenHandler(c *gin.Context)
}

type authHandler struct {
	userUseCase usecase.UserUseCase
	cache       *cache.Cache
	config      *config.Config
	logger      *logrus.Logger
	tokenMaker  token.Maker
}

func NewAuthHandler(userUsecase usecase.UserUseCase, cache *cache.Cache, config *config.Config, logger *logrus.Logger, tokenMaker token.Maker) AuthHandler {
	return &authHandler{
		userUseCase: userUsecase,
		cache:       cache,
		config:      config,
		logger:      logger,
		tokenMaker:  tokenMaker,
	}
}

func (h *authHandler) LoginHandler(c *gin.Context) {
	var loginRequest model.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	user, err := h.userUseCase.GetUser(c, 0, loginRequest.Username)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	if err := util.CheckPassword(loginRequest.Password, user.Password); err != nil {
		h.logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Username or password is incorrect"})
		return
	}

	accessToken, payload, err := h.tokenMaker.CreateToken(user.Username, string(user.Role), h.config.AccessTokenDuration)

	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	refreshToken, refreshPayload, err := h.tokenMaker.CreateToken(user.Username, string(user.Role), h.config.RefreshTokenDuration)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	session := cache.Session{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
		CreatedAt:    time.Now(),
	}

	if err := h.cache.CreateSession(session); err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}
	// send token to client through cookie
	c.Status(200)
	c.Header("Content-Type", "application/json")
	c.SetCookie("access_token", accessToken, int(h.config.AccessTokenDuration.Seconds()), "/", h.config.ServerPublicUrl, false, true)
	c.SetCookie("refresh_token", refreshToken, int(h.config.RefreshTokenDuration.Seconds()), "/", h.config.ServerPublicUrl, false, true)
	c.JSON(200, model.ResponseData[model.AuthResponse]{Code: 200, Status: "success", Data: model.AuthResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  payload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: model.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}})
}

func (h *authHandler) IsAuthHandler(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		h.logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	payload, err := h.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		h.logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}
	c.JSON(200, model.ResponseData[token.Payload]{Code: 200, Status: "success", Data: *payload})
}

func (h *authHandler) LogoutHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	refreshPayload, err := h.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		h.logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	session, err := h.cache.GetSession(refreshPayload.ID.String())

	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
	}

	if err := h.cache.DeleteSession(session.ID.String()); err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	c.SetCookie("access_token", "", -1, "/", h.config.ServerPublicUrl, false, true)
	c.SetCookie("refresh_token", "", -1, "/", h.config.ServerPublicUrl, false, true)

	c.JSON(200, model.ResponseMessage{Code: 200, Status: "success", Message: "Logout success"})
}

func (h *authHandler) RefreshAccessTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	refreshPayload, err := h.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		h.logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	session, err := h.cache.GetSession(refreshPayload.ID.String())

	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
	}

	if session.IsBlocked {
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
	}

	newAccessToken, newAccessPayload, err := h.tokenMaker.CreateToken(refreshPayload.Username, refreshPayload.Role, h.config.AccessTokenDuration)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	c.SetCookie("access_token", newAccessToken, int(h.config.AccessTokenDuration.Seconds()), "/", h.config.ServerPublicUrl, false, true)

	c.JSON(200, model.ResponseData[token.Payload]{Code: 200, Status: "success", Data: *newAccessPayload})
}
