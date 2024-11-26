package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func SantriPresenceRouting(handler handler.SantriPresenceHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/santri-presence",
			Handle:      handler.CreateSantriPresenceHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
