package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func DeviceRouter(handler handler.DeviceHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/device",
			Handle:      handler.CreateDeviceHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/device",
			Handle:      handler.ListDevicesHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/device/:id",
			Handle:      handler.UpdateDeviceHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/device/:id",
			Handle:      handler.DeleteDeviceHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
