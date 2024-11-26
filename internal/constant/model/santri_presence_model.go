package model

import (
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/go-playground/validator/v10"
)

type CreateSantriPresenceRequest struct {
	ScheduleID         int32                    `json:"schedule_id" binding:"required"`
	ScheduleName       string                   `json:"schedule_name"`
	Type               db.PresenceType          `json:"type" binding:"presencetype"`
	SantriID           int32                    `json:"santri_id" binding:"required"`
	CreatedAt          string                   `json:"created_at"`
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

func IsValidPresenceType(fl validator.FieldLevel) bool {
	presenceType := db.PresenceType(fl.Field().String())

	switch presenceType {
	case db.PresenceTypeAlpha, db.PresenceTypeLate, db.PresenceTypePermission, db.PresenceTypeSick, db.PresenceTypePresent:
		return true
	default:
		return false
	}

}
