package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
)

func SmartCardRouter(handler handler.SmartCardHandler) []routers.Route {
	return []routers.Route{
		{
			Method: http.MethodGet,
			Path:   "/smart-card",
			Handle: handler.ListSmartCardsHandler,
		},
		{
			Method: http.MethodPut,
			Path:   "/smart-card/:id",
			Handle: handler.UpdateSmartCardHandler,
		},
		{
			Method: http.MethodDelete,
			Path:   "/smart-card/:id",
			Handle: handler.DeleteSmartCardHandler,
		},
	}

}
