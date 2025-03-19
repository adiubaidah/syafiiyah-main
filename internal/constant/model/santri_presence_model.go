package model

import (
	"time"

	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/go-playground/validator/v10"
)

type CreateSantriPresenceRequest struct {
	ScheduleID         int32                      `json:"schedule_id" binding:"required"`
	ScheduleName       string                     `json:"-"`
	Type               repo.PresenceType          `json:"type" binding:"presencetype"`
	SantriID           int32                      `json:"santri_id" binding:"required"`
	Notes              string                     `json:"notes"`
	CreatedBy          repo.PresenceCreatedByType `json:"-"`
	SantriPermissionID int32                      `json:"santri_permission_id"`
}
type CreateSantriPresenceBulkRequest struct {
	ScheduleIDs []int32 `json:"schedule_ids" binding:"required"`
	Type        repo.PresenceType
	SantriID    int32 `json:"santri_id" binding:"required"`
}

type ListSantriPresenceRequest struct {
	Q          string            `form:"q"`
	Limit      int32             `form:"limit" binding:"omitempty,gte=1"`
	Page       int32             `form:"page" binding:"omitempty,gte=1"`
	ScheduleID int32             `form:"schedule"`
	SantriID   int32             `form:"santri_id"`
	Type       repo.PresenceType `form:"type" binding:"omitempty,presencetype"`
	From       string            `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To         string            `form:"to" binding:"omitempty,datetime=2006-01-02"`
}

type ListMissingSantriPresenceRequest struct {
	ScheduleID int32     `form:"schedule_id" binding:"required"`
	Time       time.Time `form:"time" binding:"required"`
}

type UpdateSantriPresenceRequest struct {
	ScheduleID         int32                      `json:"schedule_id"`
	ScheduleName       string                     `json:"-"`
	Type               repo.PresenceType          `json:"type" binding:"omitempty,presencetype"`
	SantriID           int32                      `json:"santri_id"`
	Notes              string                     `json:"notes"`
	CreatedBy          repo.PresenceCreatedByType `json:"-"`
	SantriPermissionID int32                      `json:"santri_permission_id"`
}

type SantriPresenceResponse struct {
	ID                 int32             `json:"id"`
	Type               repo.PresenceType `json:"type"`
	SantriID           int32             `json:"santri_id"`
	CreatedAt          string            `json:"created_at"`
	Notes              string            `json:"notes"`
	SantriPermissionID int32             `json:"santri_permission_id"`
	Schedule           IdAndName         `json:"schedule"`
	Santri             IdAndName         `json:"santri"`
}

type ListSantriPresenceResponse struct {
	Data       []SantriPresenceResponse `json:"data"`
	Pagination Pagination               `json:"pagination"`
}

func IsValidPresenceType(fl validator.FieldLevel) bool {
	presenceType := repo.PresenceType(fl.Field().String())

	switch presenceType {
	case repo.PresenceTypeAlpha, repo.PresenceTypeLate, repo.PresenceTypePermission, repo.PresenceTypeSick, repo.PresenceTypePresent:
		return true
	default:
		return false
	}

}
