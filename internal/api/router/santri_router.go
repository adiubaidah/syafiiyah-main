package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func SantriRouter(handler handler.SantriHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/santri",
			Handle:      handler.CreateSantriHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/santri",
			Handle:      handler.ListSantriHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/santri/:id",
			Handle:      handler.GetSantriHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/santri/:id",
			Handle:      handler.UpdateSantriHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/santri/:id",
			Handle:      handler.DeleteSantriHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
