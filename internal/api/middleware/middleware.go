package middleware

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Middleware interface {
	Auth() gin.HandlerFunc
	RequireRoles(allowedRoles ...db.RoleType) gin.HandlerFunc
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
		c.Set("user", payload)
		c.Next()
	}
}

func (m *middleware) RequireRoles(allowedRoles ...db.RoleType) gin.HandlerFunc {
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

		user, ok := userValue.(*token.Payload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ResponseMessage{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			})
			return
		}

		for _, role := range allowedRoles {
			if user.Role == string(role) {
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
