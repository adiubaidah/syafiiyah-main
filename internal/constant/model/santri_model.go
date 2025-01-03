package model

import (
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/go-playground/validator/v10"
)

type CreateSantriRequest struct {
	Nis          string          `form:"nis" binding:"required"`
	Name         string          `form:"name" binding:"required"`
	Gender       repo.GenderType `form:"gender" binding:"required,oneof=male female"`
	IsActive     string          `form:"is_active" binding:"required,oneof=true false"`
	Generation   int32           `form:"generation" binding:"required"`
	Photo        string          `form:"-"`
	OccupationID int32           `form:"occupation_id"`
	ParentID     int32           `form:"parent_id"`
}

type SantriParent struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type SantriOccupation struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type SantriCompleteResponse struct {
	ID           int32            `json:"id"`
	Nis          string           `json:"nis"`
	Name         string           `json:"name"`
	Gender       repo.GenderType  `json:"gender"`
	IsActive     bool             `json:"is_active"`
	Generation   int32            `json:"generation"`
	Photo        string           `json:"photo"`
	OccupationID int32            `json:"occupation_id"`
	ParentID     int32            `json:"parent_id"`
	Occupation   SantriOccupation `json:"occupation"`
	Parent       SantriParent     `json:"parent"`
}

type ListSantriRequest struct {
	Q            string `form:"q"`
	Order        string `form:"order" binding:"omitempty,santri-order"`
	Limit        int32  `form:"limit" binding:"omitempty,gte=1"`
	Page         int32  `form:"page" binding:"omitempty,gte=1"`
	Generation   int32  `form:"generation"`
	IsActive     int    `form:"is-active" binding:"omitempty,oneof=-1 0 1"`
	OccupationID int32  `form:"occupation_id"`
}

type UpdateSantriRequest struct {
	Nis          string          `form:"nis"`
	Name         string          `form:"name"`
	Gender       repo.GenderType `form:"gender" binding:"omitempty,oneof=male female"`
	IsActive     string          `form:"is_active" binding:"omitempty,oneof=true false"`
	Generation   int32           `form:"generation"`
	Photo        string          `form:"-"`
	OccupationID int32           `form:"occupation_id"`
	ParentID     int32           `form:"parent_id"`
}

type SantriResponse struct {
	ID           int32           `json:"id"`
	Nis          string          `json:"nis"`
	Name         string          `json:"name"`
	Gender       repo.GenderType `json:"gender"`
	IsActive     bool            `json:"is_active"`
	Generation   int32           `json:"generation"`
	Photo        string          `json:"photo"`
	OccupationID int32           `json:"occupation_id"`
	ParentID     int32           `json:"parent_id"`
}

type ListSantriResponse struct {
	Items      []SantriCompleteResponse `json:"items"`
	Pagination Pagination               `json:"pagination"`
}

func IsValidSantriOrder(fl validator.FieldLevel) bool {
	order := repo.SantriOrderBy(fl.Field().String())
	switch order {
	case repo.SantriOrderByAscName, repo.SantriOrderByDescName, repo.SantriOrderByAscGeneration, repo.SantriOrderByDescGeneration, repo.SantriOrderByAscNis, repo.SantriOrderByDescNis:
		return true
	default:
		return false
	}
}
