package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/exception"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmployeeOccuapationUsecase interface {
	CreateEmployeeOccupation(ctx context.Context, request *model.CreateEmployeeOccupationRequest) (*model.EmployeeOccupationResponse, error)
	ListEmployeeOccupations(ctx context.Context) (*[]model.EmployeeOccupationWithCountResponse, error)
	UpdateEmployeeOccupation(ctx context.Context, request *model.UpdateEmployeeOccupationRequest, occupationId int32) (*model.EmployeeOccupationResponse, error)
	DeleteEmployeeOccupation(ctx context.Context, occupationId int32) (*model.EmployeeOccupationResponse, error)
}

type employeeOccupationService struct {
	store repo.Store
}

func NewEmployeeOccupationUseCase(store repo.Store) EmployeeOccuapationUsecase {
	return &employeeOccupationService{store: store}
}

func (s *employeeOccupationService) CreateEmployeeOccupation(ctx context.Context, request *model.CreateEmployeeOccupationRequest) (*model.EmployeeOccupationResponse, error) {

	result, err := s.store.CreateEmployeeOccupation(ctx, repo.CreateEmployeeOccupationParams{
		Name:        request.Name,
		Description: pgtype.Text{String: request.Description, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &model.EmployeeOccupationResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description.String,
	}, nil
}

func (s *employeeOccupationService) ListEmployeeOccupations(ctx context.Context) (*[]model.EmployeeOccupationWithCountResponse, error) {
	results, err := s.store.ListEmployeeOccupations(ctx)
	if err != nil {
		return nil, err
	}

	var response []model.EmployeeOccupationWithCountResponse
	for _, santriOccupation := range results {
		response = append(response, model.EmployeeOccupationWithCountResponse{
			EmployeeOccupationResponse: model.EmployeeOccupationResponse{
				ID:          santriOccupation.ID,
				Name:        santriOccupation.Name,
				Description: santriOccupation.Description.String,
			},
			Count: int32(santriOccupation.Count),
		})
	}

	return &response, nil
}

func (s *employeeOccupationService) UpdateEmployeeOccupation(ctx context.Context, request *model.UpdateEmployeeOccupationRequest, occupationId int32) (*model.EmployeeOccupationResponse, error) {
	result, err := s.store.UpdateEmployeeOccupation(ctx, repo.UpdateEmployeeOccupationParams{
		ID:          occupationId,
		Name:        pgtype.Text{String: request.Name, Valid: request.Name != ""},
		Description: pgtype.Text{String: request.Description, Valid: true},
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Employee Occupation not found")
		}
		return nil, err
	}

	return &model.EmployeeOccupationResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description.String,
	}, nil
}

func (s *employeeOccupationService) DeleteEmployeeOccupation(ctx context.Context, occupationId int32) (*model.EmployeeOccupationResponse, error) {
	result, err := s.store.DeleteEmployeeOccupation(ctx, occupationId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Employee Occupation not found")
		}
		return nil, err
	}

	return &model.EmployeeOccupationResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description.String,
	}, nil
}
