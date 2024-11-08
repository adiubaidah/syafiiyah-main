package model

type BaseParentRequest struct {
	Name           string `form:"name" binding:"required"`
	Address        string `form:"address" binding:"required"`
	Gender         string `form:"gender" binding:"required,oneof=male female"`
	WhatsappNumber string `form:"whatsapp_number" binding:"min=10,max=13"`
	Photo          string `form:"photo"`
	UserID         int32  `form:"user_id" binding:"required"`
}

type ListParentRequest struct {
	Q       string `json:"q"`
	Limit   int32  `json:"limit" binding:"required"`
	Page    int32  `json:"page" binding:"required"`
	HasUser int8   `json:"has_user" binding:"oneof=0 1 -1"`
	Order   string `json:"order" binding:"oneof=asc:name desc:name"`
}
