package model

import (
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
)

type SmartCardRequest struct {
	Uid string `json:"uid" validate:"required"`
}

type ListSmartCardRequest struct {
	CardOwner repo.CardOwner `form:"card-owner" binding:"omitempty,oneof=santri employee none all"`
	IsActive  int            `form:"is-active" binding:"omitempty,oneof=-1 0 1"`
	Q         string         `form:"q"`
	Page      int32          `form:"page" binding:"omitempty,gte=1"`
	Limit     int32          `form:"limit" binding:"omitempty,gte=1"`
}

type UpdateSmartCardRequest struct {
	IsActive  bool          `json:"is_active"`
	OwnerRole repo.RoleType `json:"owner_role" binding:"omitempty,oneof=santri employee admin superadmin"` //parent can't have card
	OwnerID   int32         `json:"owner_id"`
}

type SmartCard struct {
	ID        int32  `json:"id"`
	Uid       string `json:"uid"`
	CreatedAt string `json:"create_at"`
	IsActive  bool   `json:"is_active"`
}

type ListSmartCardResponse struct {
	Items      []SmartCardComplete `json:"items"`
	Pagination Pagination          `json:"pagination"`
}

type SmartCardComplete struct {
	SmartCard
	Owner OwenerDetails `json:"owner"`
}

type OwenerDetails struct {
	ID   int32         `json:"id"`
	Role repo.RoleType `json:"role"`
	Name string        `json:"name"`
}
