package middleware

import (
	"net/http"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Middleware interface {
	Auth() gin.HandlerFunc
	RequireRoles(allowedRoles ...repo.RoleType) gin.HandlerFunc
}

type middleware struct {
	logger     *logrus.Logger
	tokenMaker token.Maker
}

func NewMiddleware(logger *logrus.Logger, tokenMaker token.Maker) Middleware {
	return &middleware{logger: logger, tokenMaker: tokenMaker}
}

func (m *middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			m.logger.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ResponseMessage{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			})
			return
		}
		payload, err := m.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			m.logger.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ResponseMessage{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			})
			return
		}
		c.Set("user", payload.User)
		c.Next()
	}
}

func (m *middleware) RequireRoles(allowedRoles ...repo.RoleType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ResponseMessage{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			})
			return
		}

		user, ok := userValue.(*model.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ResponseMessage{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			})
			return
		}

		for _, role := range allowedRoles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, model.ResponseMessage{
			Code:    http.StatusForbidden,
			Status:  "error",
			Message: "Forbidden",
		})
	}
}
