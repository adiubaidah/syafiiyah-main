package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
)

type ArduinoUseCase interface {
	CreateArduino(ctx context.Context, request *model.CreateArduinoRequest) (model.ArduinoResponse, error)
	ListArduinos(ctx context.Context) ([]model.ArduinoWithModesResponse, error)
	DeleteArduino(ctx context.Context, arduinoId int32) (model.ArduinoResponse, error)
}

type arduinoService struct {
	store db.Store
}

func NewArduinoUseCase(store db.Store) ArduinoUseCase {
	return &arduinoService{store: store}
}

func (c *arduinoService) CreateArduino(ctx context.Context, request *model.CreateArduinoRequest) (model.ArduinoResponse, error) {
	sqlStore := c.store.(*db.SQLStore)
	modeParams := make([]db.CreateArduinoModesParams, 0)

	for _, mode := range request.Modes {
		modeParams = append(modeParams, db.CreateArduinoModesParams{
			Mode:                 mode,
			InputTopic:           fmt.Sprintf("%s/input/%s", request.Name, mode),
			AcknowledgementTopic: fmt.Sprintf("%s/acknowledgment/%s", request.Name, mode),
		})
	}

	arduino, err := sqlStore.CreateArduinoWithModes(ctx, request.Name, modeParams)
	if err != nil {
		return model.ArduinoResponse{}, err
	}

	return model.ArduinoResponse{
		ID:   arduino.ID,
		Name: arduino.Name,
	}, nil
}

func (c *arduinoService) ListArduinos(ctx context.Context) ([]model.ArduinoWithModesResponse, error) {
	arduinos, err := c.store.ListArduinos(ctx)
	if err != nil {
		return nil, err
	}

	arduinoMap := make(map[int32]*model.ArduinoWithModesResponse)

	for _, arduino := range arduinos {
		if _, exists := arduinoMap[arduino.ID]; !exists {
			arduinoMap[arduino.ID] = &model.ArduinoWithModesResponse{
				ID:    arduino.ID,
				Name:  arduino.Name,
				Modes: []model.ModeArduino{},
			}
		}

		if arduino.ArduinoModeID.Valid {
			mode := model.ModeArduino{
				Mode:                arduino.ArduinoModeMode.ArduinoModeType,
				InputTopic:          arduino.ArduinoModeInputTopic.String, // Populate these fields if available
				AcknowledgmentTopic: arduino.ArduinoModeAcknowledgementTopic.String,
			}
			arduinoMap[arduino.ID].Modes = append(arduinoMap[arduino.ID].Modes, mode)
		}
	}
	var arduinoResponses []model.ArduinoWithModesResponse
	for _, arduino := range arduinoMap {
		arduinoResponses = append(arduinoResponses, *arduino)
	}

	return arduinoResponses, nil
}

func (c *arduinoService) DeleteArduino(ctx context.Context, arduinoId int32) (model.ArduinoResponse, error) {
	arduino, err := c.store.DeleteArduino(ctx, arduinoId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return model.ArduinoResponse{}, exception.NewNotFoundError("Arduino not found")
		}
		return model.ArduinoResponse{}, err
	}

	return model.ArduinoResponse{
		ID:   arduino.ID,
		Name: arduino.Name,
	}, nil
}
