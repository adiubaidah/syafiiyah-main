package model

import repo "github.com/adiubaidah/syafiiyah-main/internal/repository"

type CreateDeviceRequest struct {
	Name  string                `json:"name" binding:"required,min=3,max=50"`
	Modes []repo.DeviceModeType `json:"modes" binding:"required"`
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
	Mode                repo.DeviceModeType `json:"mode"`
	InputTopic          string              `json:"input_topic"`
	AcknowledgmentTopic string              `json:"acknowledgment_topic"`
}
