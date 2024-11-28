package model

import (
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/go-playground/validator/v10"
)

type CreateSantriPresenceRequest struct {
	ScheduleID         int32                    `json:"schedule_id" binding:"required"`
	ScheduleName       string                   `json:"-"`
	Type               db.PresenceType          `json:"type" binding:"presencetype"`
	SantriID           int32                    `json:"santri_id" binding:"required"`
	Notes              string                   `json:"notes"`
	CreatedBy          db.PresenceCreatedByType `json:"-"`
	SantriPermissionID int32                    `json:"santri_permission_id"`
}

type ListSantriPresenceRequest struct {
	Q          string          `form:"q"`
	Limit      int32           `form:"limit" binding:"omitempty,gte=1"`
	Page       int32           `form:"page" binding:"omitempty,gte=1"`
	ScheduleID int32           `form:"schedule"`
	SantriID   int32           `form:"santri_id"`
	Type       db.PresenceType `form:"type" binding:"omitempty,presencetype"`
	From       string          `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To         string          `form:"to" binding:"omitempty,datetime=2006-01-02"`
}

type UpdateSantriPresenceRequest struct {
	ScheduleID         int32                    `json:"schedule_id"`
	ScheduleName       string                   `json:"-"`
	Type               db.PresenceType          `json:"type" binding:"omitempty,presencetype"`
	SantriID           int32                    `json:"santri_id"`
	Notes              string                   `json:"notes"`
	CreatedBy          db.PresenceCreatedByType `json:"-"`
	SantriPermissionID int32                    `json:"santri_permission_id"`
}

type SantriPresenceResponse struct {
	ID                 int32           `json:"id"`
	Type               db.PresenceType `json:"type"`
	SantriID           int32           `json:"santri_id"`
	CreatedAt          string          `json:"created_at"`
	Notes              string          `json:"notes"`
	SantriPermissionID int32           `json:"santri_permission_id"`
	Schedule           IdAndName       `json:"schedule"`
	Santri             IdAndName       `json:"santri"`
}

type ListSantriPresenceResponse struct {
	Data       []SantriPresenceResponse `json:"data"`
	Pagination Pagination               `json:"pagination"`
}

func IsValidPresenceType(fl validator.FieldLevel) bool {
	presenceType := db.PresenceType(fl.Field().String())

	switch presenceType {
	case db.PresenceTypeAlpha, db.PresenceTypeLate, db.PresenceTypePermission, db.PresenceTypeSick, db.PresenceTypePresent:
		return true
	default:
		return false
	}

}
