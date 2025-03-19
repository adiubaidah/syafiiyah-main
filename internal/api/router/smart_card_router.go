package routing

import (
	"net/http"

	"github.com/adiubaidah/syafiiyah-main/internal/api/handler"
	"github.com/adiubaidah/syafiiyah-main/platform/routers"
)

func SmartCardRouter(handler *handler.SmartCardHandler) []routers.Route {
	return []routers.Route{
		{
			Method: http.MethodGet,
			Path:   "/smart-card",
			Handle: handler.List,
		},
		{
			Method: http.MethodPut,
			Path:   "/smart-card/:id",
			Handle: handler.Update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/smart-card/:id",
			Handle: handler.Delete,
		},
	}

}
