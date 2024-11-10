package model

import (
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/go-playground/validator/v10"
)

type CreateSantriRequest struct {
	Nis          string    `form:"nis" binding:"required"`
	Name         string    `form:"name" binding:"required"`
	Gender       db.Gender `form:"gender" binding:"required"`
	IsActive     string    `form:"is_active" binding:"required,oneof=true false"`
	Generation   int32     `form:"generation" binding:"required"`
	Photo        string    `form:"-"`
	OccupationID int32     `form:"occupation_id"`
	ParentID     int32     `form:"parent_id"`
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
	Gender       db.Gender        `json:"gender"`
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
	Order        string `form:"order" binding:"omitempty,santriorder"`
	Limit        int32  `form:"limit" binding:"omitempty,gte=1"`
	Page         int32  `form:"page" binding:"omitempty,gte=1"`
	Generation   int32  `form:"generation"`
	IsActive     int    `form:"is-active"`
	OccupationID int32  `form:"occupation_id"`
}

type UpdateSantriRequest struct {
	CreateSantriRequest
}

type SantriResponse struct {
	ID           int32     `json:"id"`
	Nis          string    `json:"nis"`
	Name         string    `json:"name"`
	Gender       db.Gender `json:"gender"`
	IsActive     bool      `json:"is_active"`
	Generation   int32     `json:"generation"`
	Photo        string    `json:"photo"`
	OccupationID int32     `json:"occupation_id"`
	ParentID     int32     `json:"parent_id"`
}

type ListSantriResponse struct {
	Items      []SantriCompleteResponse `json:"items"`
	Pagination Pagination               `json:"pagination"`
}

func IsValidSantriOrder(fl validator.FieldLevel) bool {
	order := db.SantriOrderBy(fl.Field().String())
	switch order {
	case db.SantriOrderByAscName, db.SantriOrderByDescName, db.SantriOrderByAscGeneration, db.SantriOrderByDescGeneration, db.SantriOrderByAscNis, db.SantriOrderByDescNis:
		return true
	default:
		return false
	}
}
