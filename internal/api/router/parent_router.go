package routing

import (
	"net/http"

	"github.com/adiubaidah/rfid-syafiiyah/internal/api/handler"
	"github.com/adiubaidah/rfid-syafiiyah/internal/api/middleware"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/adiubaidah/rfid-syafiiyah/platform/routers"
	"github.com/gin-gonic/gin"
)

func ParentRouter(middle middleware.Middleware, handler *handler.ParentHandler) []routers.Route {
	return []routers.Route{
		{
			Method: http.MethodPost,
			Path:   "/parent",
			Handle: handler.CreateParentHandler,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/parent",
			Handle: handler.ListParentHandler,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/parent/:id",
			Handle: handler.GetParentHandler,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
		{
			Method: http.MethodPut,
			Path:   "/parent/:id",
			Handle: handler.UpdateParentHandler,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
		{
			Method: http.MethodDelete,
			Path:   "/parent/:id",
			Handle: handler.DeleteParentHandler,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
	}
}
