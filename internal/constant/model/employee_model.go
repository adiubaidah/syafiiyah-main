package model

import db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"

type CreateEmployeeRequest struct {
	Name   string        `form:"name" binding:"required"`
	NIP    string        `form:"nip" binding:"required,numeric,len=18"`
	Gender db.GenderType `form:"gender" binding:"required,oneof=male female"`
}
