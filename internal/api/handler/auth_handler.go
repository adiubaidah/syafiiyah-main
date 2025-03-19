package handler

import (
	"time"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/exception"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	"github.com/adiubaidah/syafiiyah-main/internal/usecase"
	"github.com/adiubaidah/syafiiyah-main/pkg/config"
	"github.com/adiubaidah/syafiiyah-main/pkg/token"
	"github.com/adiubaidah/syafiiyah-main/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
)

type AuthHandler struct {
	Config         *config.Config
	UserUseCase    *usecase.UserUseCase
	SessionUseCase *usecase.SessionUseCase
	Logger         *logrus.Logger
	TokenMaker     token.Maker
}

func NewAuthHandler(args *AuthHandler) *AuthHandler {
	return args
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var request model.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	if err := request.Validate(); err != nil {
		h.Logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	var user *model.User

	if request.Username != "" {
		result, err := h.UserUseCase.GetByUsername(c, request.Username)
		if err != nil {
			h.Logger.Error(err)
			if appErr, ok := err.(*exception.AppError); ok {
				c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
				return
			}
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
			return
		}
		if err := util.CheckPassword(request.Password, result.Password); err != nil {
			h.Logger.Error(err)
			c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Username or password is incorrect"})
			return
		}

		user = &model.User{
			ID:       result.ID,
			Username: result.Username,
			Role:     result.Role,
		}

	} else if request.Token != "" {
		payload, err := idtoken.Validate(c, request.Token, h.Config.GoogleOauthClient)
		if err != nil {
			h.Logger.Error(err)
			c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
			return
		}
		result, err := h.UserUseCase.GetByEmail(c, payload.Claims["email"].(string))

		if err != nil {
			h.Logger.Error(err)
			if appErr, ok := err.(*exception.AppError); ok {
				c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
				return
			}
			c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
			return
		}

		user = &model.User{
			ID:       result.ID,
			Username: result.Username,
			Role:     result.Role,
		}
	} else {
		h.Logger.Error("Invalid login request")
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: "Invalid login request"})
		return
	}

	accessToken, payload, err := h.TokenMaker.CreateToken(user, h.Config.AccessTokenDuration)

	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	refreshToken, refreshPayload, err := h.TokenMaker.CreateToken(user, h.Config.RefreshTokenDuration)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	session := model.Session{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
		CreatedAt:    time.Now(),
	}

	if err := h.SessionUseCase.Create(session); err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}
	// send token to client through cookie
	c.Status(200)
	c.Header("Content-Type", "application/json")
	c.SetCookie("access_token", accessToken, int(h.Config.AccessTokenDuration.Seconds()), "/", h.Config.ServerPublicUrl, false, true)
	c.SetCookie("refresh_token", refreshToken, int(h.Config.RefreshTokenDuration.Seconds()), "/", h.Config.ServerPublicUrl, false, true)
	c.JSON(200, model.ResponseData[model.AuthResponse]{Code: 200, Status: "success", Data: model.AuthResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  payload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: model.User{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}})
}

func (h *AuthHandler) IsAuthHandler(c *gin.Context) {
	userValue, _ := c.Get("user")
	user, _ := userValue.(*model.User)
	c.JSON(200, model.ResponseData[*model.User]{Code: 200, Status: "success", Data: user})
}

func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.Logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	refreshPayload, err := h.TokenMaker.VerifyToken(refreshToken)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	session, err := h.SessionUseCase.Get(refreshPayload.ID.String())

	if err != nil {
		h.Logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
	}

	if err := h.SessionUseCase.Delete(session.ID.String()); err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	c.SetCookie("access_token", "", -1, "/", h.Config.ServerPublicUrl, false, true)
	c.SetCookie("refresh_token", "", -1, "/", h.Config.ServerPublicUrl, false, true)

	c.JSON(200, model.ResponseMessage{Code: 200, Status: "success", Message: "Logout success"})
}

func (h *AuthHandler) RefreshAccessTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.Logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	refreshPayload, err := h.TokenMaker.VerifyToken(refreshToken)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(401, model.ResponseMessage{Code: 401, Status: "error", Message: "Unauthorized"})
		return
	}

	session, err := h.SessionUseCase.Get(refreshPayload.ID.String())

	if err != nil {
		h.Logger.Error(err)
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

	newAccessToken, newAccessPayload, err := h.TokenMaker.CreateToken(&model.User{
		ID:       refreshPayload.User.ID,
		Username: refreshPayload.User.Username,
		Role:     refreshPayload.User.Role,
	}, h.Config.AccessTokenDuration)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	c.SetCookie("access_token", newAccessToken, int(h.Config.AccessTokenDuration.Seconds()), "/", h.Config.ServerPublicUrl, false, true)

	c.JSON(200, model.ResponseData[token.Payload]{Code: 200, Status: "success", Data: *newAccessPayload})
}
