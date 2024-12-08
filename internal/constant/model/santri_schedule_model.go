package model

import (
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/go-playground/validator/v10"
)

type CreateSantriScheduleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	StartPresence string `json:"start_presence" binding:"validTime"`
	StartTime     string `json:"start_time" binding:"validTime"`
	FinishTime    string `json:"finish_time" binding:"validTime"`
}

type UpdateSantriScheduleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	StartPresence string `json:"start_presence" binding:"omitempty,validTime"`
	StartTime     string `json:"start_time" binding:"omitempty,validTime"`
	FinishTime    string `json:"finish_time" binding:"omitempty,validTime"`
}

func IsValidTime(fl validator.FieldLevel) bool {
	_, err := util.ParseTime(fl.Field().String())
	return err == nil
}

type SantriScheduleResponse struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	StartPresence string `json:"start_presence"`
	StartTime     string `json:"start_time"`
	FinishTime    string `json:"finish_time"`
}
