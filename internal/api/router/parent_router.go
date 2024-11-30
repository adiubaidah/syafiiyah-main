package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func ParentRouter(handler handler.ParentHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/parent",
			Handle:      handler.CreateParentHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/parent",
			Handle:      handler.ListParentHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/parent/:id",
			Handle:      handler.GetParentHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/parent/:id",
			Handle:      handler.UpdateParentHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/parent/:id",
			Handle:      handler.DeleteParentHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
