package model

type BaseSantriOccupation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateSantriOccupationRequest struct {
	BaseSantriOccupation
}

type UpdateSantriOccupationRequest struct {
	BaseSantriOccupation
	ID int32 `json:"-" binding:"required"`
}

type SantriOccupationResponse struct {
	BaseSantriOccupation
	ID int32 `json:"id"`
}

type SantriOccupationWithCountResponse struct {
	BaseSantriOccupation
	Count int32 `json:"count"`
}

type ListSantriOccupationParams struct {
	Q string `json:"q"`
}
