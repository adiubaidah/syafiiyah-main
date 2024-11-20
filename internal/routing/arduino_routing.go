package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func ArduinoRouting(handler handler.ArduinoHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/arduino",
			Handle:      handler.CreateArduinoHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/arduino",
			Handle:      handler.ListArduinosHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/arduino/:id",
			Handle:      handler.UpdateArduinoHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/arduino/:id",
			Handle:      handler.DeleteArduinoHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
