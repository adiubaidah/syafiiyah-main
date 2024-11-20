package model

type CreateSantriOccupationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateSantriOccupationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SantriOccupationResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SantriOccupationWithCountResponse struct {
	SantriOccupationResponse
	Count int32 `json:"count"`
}
