package model

import db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"

type CreateArduinoRequest struct {
	Name  string               `json:"name" binding:"required,min=3,max=50"`
	Modes []db.ArduinoModeType `json:"modes" binding:"required"`
}

type ArduinoResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type ArduinoWithModesResponse struct {
	ID    int32         `json:"id"`
	Name  string        `json:"name"`
	Modes []ModeArduino `json:"modes"`
}

type ModeArduino struct {
	Mode                db.ArduinoModeType `json:"mode"`
	InputTopic          string             `json:"input_topic"`
	AcknowledgmentTopic string             `json:"acknowledgment_topic"`
}
