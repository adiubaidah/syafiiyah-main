package routing

import (
	"net/http"

	"github.com/adiubaidah/syafiiyah-main/internal/api/handler"
	"github.com/adiubaidah/syafiiyah-main/internal/api/middleware"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/platform/routers"
	"github.com/gin-gonic/gin"
)

func EmployeeRouter(middle middleware.Middleware, handler *handler.EmployeeHandler) []routers.Route {
	return []routers.Route{
		{
			Method: http.MethodPost,
			Path:   "/employee",
			Handle: handler.Create,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/employee",
			Handle: handler.List,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
		{
			Method: http.MethodPut,
			Path:   "/employee/:id",
			Handle: handler.Update,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
		{
			Method: http.MethodDelete,
			Path:   "/employee/:id",
			Handle: handler.Delete,
			MiddleWares: []gin.HandlerFunc{
				middle.Auth(),
				middle.RequireRoles(repo.RoleTypeAdmin, repo.RoleTypeSuperadmin),
			},
		},
	}
}
