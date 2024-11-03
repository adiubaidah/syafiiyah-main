package model

type CreateSantriOccupationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateSantriOccupationRequest struct {
	ID          int32  `json:"-" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type SantriOccupationResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SantriOccupationWithCountResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int32  `json:"count"`
}

type ListSantriOccupationParams struct {
	Q string `json:"q"`
}
