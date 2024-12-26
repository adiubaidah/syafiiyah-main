package model

import (
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Username string        `json:"username" binding:"required"`
	Role     repo.RoleType `json:"role" binding:"required,role"`
	Password string        `json:"password" binding:"required"`
}

type User struct {
	ID       int32         `json:"id"`
	Username string        `json:"username"`
	Role     repo.RoleType `json:"role"`
}

type UserWithPassword struct {
	ID       int32
	Username string
	Role     repo.RoleType
	Password string
}

type ListUserRequest struct {
	Q        string        `form:"q"`
	Order    string        `form:"order" binding:"omitempty,userorder"`
	Limit    int32         `form:"limit" binding:"omitempty,gte=1"`
	Page     int32         `form:"page" binding:"omitempty,gte=1"`
	HasOwner int32         `form:"has-owner"`
	Role     repo.RoleType `form:"role" binding:"omitempty,role"`
}

type UpdateUserRequest struct {
	Username string        `json:"username"`
	Role     repo.RoleType `json:"role" binding:"omitempty,role"`
	Password string        `json:"password"`
}

type UserComplete struct {
	ID          int32  `json:"id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	UserDetails `json:"details"`
}

type UserDetails struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type ListUserResponse struct {
	Items      []UserComplete `json:"items"`
	Pagination Pagination     `json:"pagination"`
}

func IsValidUserOrder(fl validator.FieldLevel) bool {
	order := repo.UserOrderBy(fl.Field().String())
	switch order {
	case repo.UserOrderByAscUsername, repo.UserOrderByDescUsername, repo.UserOrderByAscName, repo.UserOrderByDescName:
		return true
	default:
		return false
	}
}

func IsValidRole(fl validator.FieldLevel) bool {
	role := repo.RoleType(fl.Field().String())
	switch role {
	case repo.RoleTypeAdmin, repo.RoleTypeEmployee, repo.RoleTypeParent, repo.RoleTypeSuperadmin, repo.RoleTypeSantri:
		return true
	default:
		return false
	}
}
