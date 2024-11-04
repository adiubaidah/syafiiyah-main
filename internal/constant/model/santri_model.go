package model

type CreateSantriRequest struct {
	Nis          string `form:"nis" binding:"required"`
	Name         string `form:"name" binding:"required"`
	Gender       string `form:"gender" binding:"required"`
	IsActive     bool   `form:"is_active" `
	Generation   int32  `form:"generation" binding:"required"`
	Photo        string `form:"photo"`
	OccupationID int32  `form:"occupation_id" binding:"required"`
	ParentID     int32  `form:"parent_id" binding:"required"`
}
