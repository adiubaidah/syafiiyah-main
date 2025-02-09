package model

import (
	"time"

	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
)

type CreateEmployeePresenceRequest struct {
	ScheduleID           int32                      `json:"schedule_id" binding:"required"`
	ScheduleName         string                     `json:"-"`
	Type                 repo.PresenceType          `json:"type" binding:"presencetype"`
	EmployeeID           int32                      `json:"employee_id" binding:"required"`
	Notes                string                     `json:"notes"`
	CreatedBy            repo.PresenceCreatedByType `json:"-"`
	EmployeePermissionID int32                      `json:"employee_permission_id"`
}
type CreateEmployeePresenceBulkRequest struct {
	ScheduleIDs []int32 `json:"schedule_ids" binding:"required"`
	Type        repo.PresenceType
	EmployeeID  int32 `json:"employee_id" binding:"required"`
}

type ListEmployeePresenceRequest struct {
	Q          string            `form:"q"`
	Limit      int32             `form:"limit" binding:"omitempty,gte=1"`
	Page       int32             `form:"page" binding:"omitempty,gte=1"`
	ScheduleID int32             `form:"schedule"`
	EmployeeID int32             `form:"employee_id"`
	Type       repo.PresenceType `form:"type" binding:"omitempty,presencetype"`
	From       string            `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To         string            `form:"to" binding:"omitempty,datetime=2006-01-02"`
}

type ListMissingEmployeePresenceRequest struct {
	ScheduleID int32     `form:"schedule_id" binding:"required"`
	Time       time.Time `form:"time" binding:"required"`
}

type UpdateEmployeePresenceRequest struct {
	ScheduleID           int32                      `json:"schedule_id"`
	ScheduleName         string                     `json:"-"`
	Type                 repo.PresenceType          `json:"type" binding:"omitempty,presencetype"`
	EmployeeID           int32                      `json:"employee_id"`
	Notes                string                     `json:"notes"`
	CreatedBy            repo.PresenceCreatedByType `json:"-"`
	EmployeePermissionID int32                      `json:"employee_permission_id"`
}

type EmployeePresenceResponse struct {
	ID                   int32             `json:"id"`
	Type                 repo.PresenceType `json:"type"`
	EmployeeID           int32             `json:"employee_id"`
	CreatedAt            string            `json:"created_at"`
	Notes                string            `json:"notes"`
	EmployeePermissionID int32             `json:"employee_permission_id"`
	Schedule             IdAndName         `json:"schedule"`
	Employee             IdAndName         `json:"Employee"`
}

type ListEmployeePresenceResponse struct {
	Data       []SantriPresenceResponse `json:"data"`
	Pagination Pagination               `json:"pagination"`
}
