package model

type CreateEmployeeOccupationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateEmployeeOccupationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EmployeeOccupationResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EmployeeOccupationWithCountResponse struct {
	EmployeeOccupationResponse
	Count int32 `json:"count"`
}
