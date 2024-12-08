package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/internal/api/middleware"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func AuthRouter(middle middleware.Middleware, handler handler.AuthHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/auth/login",
			Handle:      handler.LoginHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method: http.MethodPost,
			Path:   "/auth/is-auth",
			Handle: handler.IsAuthHandler,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/auth/logout",
			Handle:      handler.LogoutHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/auth/refresh-access-token",
			Handle:      handler.RefreshAccessTokenHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
