package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func HolidayRouter(handler handler.HolidayHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/holiday",
			Handle:      handler.CreateHolidayHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/holiday",
			Handle:      handler.ListHolidaysHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/holiday/:id",
			Handle:      handler.UpdateHolidayHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/holiday/:id",
			Handle:      handler.DeleteHolidayHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
