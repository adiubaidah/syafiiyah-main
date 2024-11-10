package model

import (
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"userrole"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type ListUserRequest struct {
	Q        string `form:"q"`
	Order    string `form:"order" binding:"omitempty,userorder"`
	Limit    int32  `form:"limit" binding:"omitempty,gte=1"`
	Page     int32  `form:"page" binding:"omitempty,gte=1"`
	HasOwner int32  `form:"has-owner"`
	Role     string `form:"role" binding:"omitempty,userrole"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Role     string `json:"role" binding:"userrole"`
	Password string `json:"password"`
}

type UserComplete struct {
	ID          int32  `json:"id"`
	Username    string `json:"username"`
	Role        string `json:"role" binding:"userrole"`
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

func IsValidUserRole(fl validator.FieldLevel) bool {
	role := db.UserRole(fl.Field().String())
	switch role {
	case db.UserRoleSuperadmin, db.UserRoleAdmin, db.UserRoleEmployee, db.UserRoleParent:
		return true
	default:
		return false
	}
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
