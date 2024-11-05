package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func SantriOccupationRouting(handler handler.SantriOccupationHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/santri-occupation",
			Handle:      handler.CreateSantriOccupationHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
