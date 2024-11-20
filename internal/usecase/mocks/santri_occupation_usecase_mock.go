package mocks

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/stretchr/testify/mock"
)

type SantriOccuapationUsecase struct {
	mock.Mock
}

func (m *SantriOccuapationUsecase) CreateSantriOccupation(ctx context.Context, request *model.CreateSantriOccupationRequest) (model.SantriOccupationResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(model.SantriOccupationResponse), args.Error(1)
}

func (m *SantriOccuapationUsecase) ListSantriOccupations(ctx context.Context) ([]model.SantriOccupationWithCountResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.SantriOccupationWithCountResponse), args.Error(1)
}

func (m *SantriOccuapationUsecase) UpdateSantriOccupation(ctx context.Context, request *model.UpdateSantriOccupationRequest, santriOccupationid int32) (model.SantriOccupationResponse, error) {
	args := m.Called(ctx, request, santriOccupationid)
	return args.Get(0).(model.SantriOccupationResponse), args.Error(1)
}

func (m *SantriOccuapationUsecase) DeleteSantriOccupation(ctx context.Context, santriOccupationId int32) (model.SantriOccupationResponse, error) {
	args := m.Called(ctx, santriOccupationId)
	return args.Get(0).(model.SantriOccupationResponse), args.Error(1)
}
