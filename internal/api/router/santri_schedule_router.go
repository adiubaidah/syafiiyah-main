package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func SantriScheduleRouter(handler handler.SantriScheduleHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/santri-schedule",
			Handle:      handler.CreateSantriScheduleHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/santri-schedule",
			Handle:      handler.ListSantriScheduleHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/santri-schedule/:id",
			Handle:      handler.UpdateSantriScheduleHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/santri-schedule/:id",
			Handle:      handler.DeleteSantriScheduleHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
