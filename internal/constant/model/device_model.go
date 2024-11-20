package model

import db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"

type CreateDeviceRequest struct {
	Name  string              `json:"name" binding:"required,min=3,max=50"`
	Modes []db.DeviceModeType `json:"modes" binding:"required"`
}

type UpdateDeviceRequest = CreateDeviceRequest

type DeviceResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type DeviceWithModesResponse struct {
	ID    int32        `json:"id"`
	Name  string       `json:"name"`
	Modes []DeviceMode `json:"modes"`
}

type DeviceMode struct {
	Mode                db.DeviceModeType `json:"mode"`
	InputTopic          string            `json:"input_topic"`
	AcknowledgmentTopic string            `json:"acknowledgment_topic"`
}
