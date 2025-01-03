package model

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CreateSantriScheduleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	StartPresence string `json:"start_presence" binding:"valid-time"`
	StartTime     string `json:"start_time" binding:"valid-time"`
	FinishTime    string `json:"finish_time" binding:"valid-time"`
}

type UpdateSantriScheduleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	StartPresence string `json:"start_presence" binding:"omitempty,valid-time"`
	StartTime     string `json:"start_time" binding:"omitempty,valid-time"`
	FinishTime    string `json:"finish_time" binding:"omitempty,valid-time"`
}

func IsValidTime(fl validator.FieldLevel) bool {
	pattern := `^\d{2}:\d{2}$`
	matched, _ := regexp.MatchString(pattern, fl.Field().String())
	return matched
}

type SantriScheduleResponse struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	StartPresence string `json:"start_presence"`
	StartTime     string `json:"start_time"`
	FinishTime    string `json:"finish_time"`
}
