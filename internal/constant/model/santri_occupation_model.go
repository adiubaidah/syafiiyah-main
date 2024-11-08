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
}

type SantriOccupationResponse struct {
	ID int32 `json:"id"`
	BaseSantriOccupation
}

type SantriOccupationWithCountResponse struct {
	ID int32 `json:"id"`
	BaseSantriOccupation
	Count int32 `json:"count"`
}
