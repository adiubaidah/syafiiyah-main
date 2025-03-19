package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/exception"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type SantriOccuapationUsecase interface {
	CreateSantriOccupation(ctx context.Context, request *model.CreateSantriOccupationRequest) (*model.SantriOccupationResponse, error)
	ListSantriOccupations(ctx context.Context) (*[]model.SantriOccupationWithCountResponse, error)
	UpdateSantriOccupation(ctx context.Context, request *model.UpdateSantriOccupationRequest, occupationId int32) (*model.SantriOccupationResponse, error)
	DeleteSantriOccupation(ctx context.Context, occupationId int32) (*model.SantriOccupationResponse, error)
}

type santriOccupationService struct {
	store repo.Store
}

func NewSantriOccupationUseCase(store repo.Store) SantriOccuapationUsecase {
	return &santriOccupationService{store: store}
}

func (s *santriOccupationService) CreateSantriOccupation(ctx context.Context, request *model.CreateSantriOccupationRequest) (*model.SantriOccupationResponse, error) {

	result, err := s.store.CreateSantriOccupation(ctx, repo.CreateSantriOccupationParams{
		Name:        request.Name,
		Description: pgtype.Text{String: request.Description, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &model.SantriOccupationResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description.String,
	}, nil
}

func (s *santriOccupationService) ListSantriOccupations(ctx context.Context) (*[]model.SantriOccupationWithCountResponse, error) {
	result, err := s.store.ListSantriOccupations(ctx)
	if err != nil {
		return nil, err
	}

	var response []model.SantriOccupationWithCountResponse
	for _, santriOccupation := range result {
		response = append(response, model.SantriOccupationWithCountResponse{
			SantriOccupationResponse: model.SantriOccupationResponse{
				ID:          santriOccupation.ID,
				Name:        santriOccupation.Name,
				Description: santriOccupation.Description.String,
			},
			Count: int32(santriOccupation.Count),
		})
	}

	return &response, nil
}

func (s *santriOccupationService) UpdateSantriOccupation(ctx context.Context, request *model.UpdateSantriOccupationRequest, occupationId int32) (*model.SantriOccupationResponse, error) {
	result, err := s.store.UpdateSantriOccupation(ctx, repo.UpdateSantriOccupationParams{
		ID:          occupationId,
		Name:        pgtype.Text{String: request.Name, Valid: request.Name != ""},
		Description: pgtype.Text{String: request.Description, Valid: true},
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Occupation not found")
		}
		return nil, err
	}

	return &model.SantriOccupationResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description.String,
	}, nil
}

func (s *santriOccupationService) DeleteSantriOccupation(ctx context.Context, occupationId int32) (*model.SantriOccupationResponse, error) {
	result, err := s.store.DeleteSantriOccupation(ctx, occupationId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Occupation not found")
		}
		return nil, err
	}

	return &model.SantriOccupationResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description.String,
	}, nil
}
