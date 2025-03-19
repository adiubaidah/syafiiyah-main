package routing

import (
	"net/http"

	"github.com/adiubaidah/syafiiyah-main/internal/api/handler"
	"github.com/adiubaidah/syafiiyah-main/platform/routers"
	"github.com/gin-gonic/gin"
)

func SantriPresenceRouter(handler handler.SantriPresenceHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/santri-presence",
			Handle:      handler.CreateSantriPresenceHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/santri-presence",
			Handle:      handler.ListSantriPresencesHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/santri-presence/:id",
			Handle:      handler.UpdateSantriPresenceHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/santri-presence/:id",
			Handle:      handler.DeleteSantriPresenceHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
