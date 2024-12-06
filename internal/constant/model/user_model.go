package model

import (
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Username string      `json:"username" binding:"required"`
	Role     db.RoleType `json:"role" binding:"required,role"`
	Password string      `json:"password" binding:"required"`
}

type User struct {
	ID       int32       `json:"id"`
	Username string      `json:"username"`
	Role     db.RoleType `json:"role"`
}

type UserWithPassword struct {
	ID       int32
	Username string
	Role     db.RoleType
	Password string
}

type ListUserRequest struct {
	Q        string      `form:"q"`
	Order    string      `form:"order" binding:"omitempty,userorder"`
	Limit    int32       `form:"limit" binding:"omitempty,gte=1"`
	Page     int32       `form:"page" binding:"omitempty,gte=1"`
	HasOwner int32       `form:"has-owner"`
	Role     db.RoleType `form:"role" binding:"omitempty,role"`
}

type UpdateUserRequest struct {
	Username string      `json:"username"`
	Role     db.RoleType `json:"role" binding:"omitempty,role"`
	Password string      `json:"password"`
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
	order := db.UserOrderBy(fl.Field().String())
	switch order {
	case db.UserOrderByAscUsername, db.UserOrderByDescUsername, db.UserOrderByAscName, db.UserOrderByDescName:
		return true
	default:
		return false
	}
}

func IsValidRole(fl validator.FieldLevel) bool {
	role := db.RoleType(fl.Field().String())
	switch role {
	case db.RoleTypeAdmin, db.RoleTypeEmployee, db.RoleTypeParent, db.RoleTypeSuperadmin, db.RoleTypeSantri:
		return true
	default:
		return false
	}
}
