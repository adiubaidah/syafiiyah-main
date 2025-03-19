package routing

import (
	"net/http"

	"github.com/adiubaidah/syafiiyah-main/internal/api/handler"
	"github.com/adiubaidah/syafiiyah-main/platform/routers"
	"github.com/gin-gonic/gin"
)

func UserRouter(handler *handler.UserHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/user",
			Handle:      handler.Create,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/user",
			Handle:      handler.List,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/user/:id",
			Handle:      handler.GetUserHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/user/:id",
			Handle:      handler.UpdateUserHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/user/:id",
			Handle:      handler.DeleteUserHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
