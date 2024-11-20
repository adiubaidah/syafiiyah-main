package model

import (
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/go-playground/validator/v10"
)

type CreateParentRequest struct {
	Name           string    `form:"name" binding:"required"`
	Address        string    `form:"address" binding:"required"`
	Gender         db.Gender `form:"gender" binding:"required,oneof=male female"`
	WhatsappNumber string    `form:"whatsapp_number" binding:"min=10,max=13"`
	Photo          string    `form:"-"`
	UserID         int32     `form:"user_id"`
}

type ListParentRequest struct {
	Q       string `form:"q"`
	Limit   int32  `form:"limit" binding:"omitempty,gte=1"`
	Page    int32  `form:"page" binding:"omitempty,gte=1"`
	HasUser int8   `form:"has-user" binding:"oneof=0 1 -1"`
	Order   string `form:"order" binding:"omitempty,parentorder"`
}

type UpdateParentRequest = CreateParentRequest

type ParentResponse struct {
	ID             int32  `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	Gender         string `json:"gender"`
	WhatsappNumber string `json:"whatsapp_number"`
	Photo          string `json:"photo"`
	UserID         int32  `json:"user_id"`
}
type ParentUser struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}

type ParentCompleteResponse struct {
	ParentResponse
	User ParentUser `json:"user"`
}

type ListParentResponse struct {
	Items      []ParentCompleteResponse `json:"items"`
	Pagination Pagination               `json:"pagination"`
}

func IsValidParentOrder(fl validator.FieldLevel) bool {
	order := db.ParentOrderBy(fl.Field().String())
	switch order {
	case db.ParentOrderByAscName, db.ParentOrderByDescName:
		return true
	default:
		return false
	}
}
