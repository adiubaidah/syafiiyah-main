package model

type CreateHolidayRequest struct {
	Name        string   `json:"name" binding:"required"`
	Color       string   `json:"color" binding:"omitempty,hexcolor"`
	Description string   `json:"description"`
	Dates       []string `json:"dates" binding:"required,dive,datetime=2006-01-02"`
}

type UpdateHolidayRequest = CreateHolidayRequest

type ListHolidayRequest struct {
	Month int32 `form:"month" binding:"omitempty,min=1,max=12"`
	Year  int32 `form:"year" binding:"omitempty,min=1"`
}

type HolidayResponse struct {
	ID          int32    `json:"id"`
	Name        string   `json:"name"`
	Color       string   `json:"color"`
	Description string   `json:"description"`
	Dates       []string `json:"dates"`
}
