package model

import db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"

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
	Q            string `json:"q"`
	Order        string `json:"order" binding:"oneof=asc:name desc:name asc:generation desc:generation asc:occupation desc:occupation asc:nis desc:nis"`
	Limit        int32  `json:"limit" binding:"required"`
	Page         int32  `json:"page" binding:"required"`
	Generation   int32  `json:"generation"`
	IsActive     int    `json:"is_active" binding:"oneof=0 1 -1"`
	OccupationID int32  `json:"occupation"`
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
