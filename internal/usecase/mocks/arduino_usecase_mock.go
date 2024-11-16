package mocks

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/stretchr/testify/mock"
)

type ArduinoUseCase struct {
	mock.Mock
}

func (m *ArduinoUseCase) CreateArduino(ctx context.Context, request *model.CreateArduinoRequest) (model.ArduinoResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(model.ArduinoResponse), args.Error(1)
}

func (m *ArduinoUseCase) ListArduinos(ctx context.Context) ([]model.ArduinoWithModesResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.ArduinoWithModesResponse), args.Error(1)
}

func (m *ArduinoUseCase) DeleteArduino(ctx context.Context, arduinoId int32) (model.ArduinoResponse, error) {

	args := m.Called(ctx, arduinoId)

	return args.Get(0).(model.ArduinoResponse), args.Error(1)

}
