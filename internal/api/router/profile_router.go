package routing

import (
	"net/http"

	"github.com/adiubaidah/syafiiyah-main/internal/api/handler"
	"github.com/adiubaidah/syafiiyah-main/internal/api/middleware"
	"github.com/adiubaidah/syafiiyah-main/platform/routers"
	"github.com/gin-gonic/gin"
)

func ProfileRouter(middle middleware.Middleware, handler *handler.ProfileHandler) []routers.Route {
	return []routers.Route{
		{
			Method: http.MethodGet,
			Path:   "/profile",
			Handle: handler.GetProfile,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
			},
		},
	}
}
