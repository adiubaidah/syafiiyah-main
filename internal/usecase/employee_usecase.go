package usecase

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmployeeUseCase struct {
	store repo.Store
}

func NewEmployeeUseCase(store repo.Store) *EmployeeUseCase {
	return &EmployeeUseCase{
		store: store,
	}
}

func (s *EmployeeUseCase) Create(ctx context.Context, request *model.CreateEmployeeRequest) (*model.Employee, error) {
	result, err := s.store.CreateEmployee(ctx, repo.CreateEmployeeParams{
		Nip:          pgtype.Text{String: request.NIP, Valid: request.NIP != ""},
		Name:         request.Name,
		Gender:       request.Gender,
		Photo:        pgtype.Text{String: request.Photo, Valid: request.Photo != ""},
		OccupationID: request.OccupationID,
		UserID:       pgtype.Int4{Int32: request.UserID, Valid: request.UserID != 0},
	})

	if err != nil {
		return nil, err
	}

	return &model.Employee{
		ID:           result.ID,
		Name:         result.Name,
		NIP:          result.Nip.String,
		Gender:       request.Gender,
		Photo:        result.Photo.String,
		OccupationID: result.OccupationID,
		UserID:       result.UserID.Int32,
	}, nil
}

func (s *EmployeeUseCase) List(ctx context.Context, request *model.ListEmployeeRequest) (*[]model.EmployeeComplete, error) {

	arg := repo.ListEmployeesParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		HasUser:      pgtype.Bool{Bool: request.HasUser == 1, Valid: request.HasUser != 0},
		LimitNumber:  request.Limit,
		OffsetNumber: (request.Page - 1) * request.Limit,
		OrderBy:      repo.NullEmployeeOrderBy{EmployeeOrderBy: repo.EmployeeOrderBy(request.Order), Valid: request.Order != ""},
	}

	employees, err := s.store.ListEmployees(ctx, arg)
	if err != nil {
		return nil, err
	}

	var result []model.EmployeeComplete
	for _, employee := range employees {
		result = append(result, model.EmployeeComplete{
			ID:           employee.ID,
			Name:         employee.Name,
			NIP:          employee.Nip.String,
			Gender:       employee.Gender,
			Photo:        employee.Photo.String,
			OccupationID: employee.OccupationID,
			UserID:       employee.UserID.Int32,
			Occupation: model.EmployeeOccupation{
				ID:   employee.OccupationID,
				Name: employee.OccupationName,
			},
		})
	}

	return &result, nil

}

func (s *EmployeeUseCase) GetByID(ctx context.Context, employeeId int32) (*model.Employee, error) {
	employee, err := s.store.GetEmployeeByID(ctx, employeeId)
	if err != nil {
		return nil, err
	}

	return &model.Employee{
		ID:   employee.ID,
		Name: employee.Name,
		NIP:  employee.Nip.String,
	}, nil
}

func (s *EmployeeUseCase) GetByUserID(ctx context.Context, userId int32) (*model.Employee, error) {
	employee, err := s.store.GetEmployeeByUserID(ctx, pgtype.Int4{Int32: userId, Valid: true})
	if err != nil {
		return nil, err
	}

	return &model.Employee{
		ID:   employee.ID,
		Name: employee.Name,
		NIP:  employee.Nip.String,
	}, nil
}

func (s *EmployeeUseCase) Count(ctx context.Context, request *model.ListEmployeeRequest) (int64, error) {
	count, err := s.store.CountEmployees(ctx, repo.CountEmployeesParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		HasUser:      pgtype.Bool{Bool: request.HasUser == 1, Valid: request.HasUser != 0},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *EmployeeUseCase) Update(ctx context.Context, request *model.UpdateEmployeeRequest, employeeId int32) (*model.Employee, error) {
	result, err := s.store.UpdateEmployee(ctx, repo.UpdateEmployeeParams{
		ID:           employeeId,
		Nip:          pgtype.Text{String: request.NIP, Valid: true}, // Nip can be null
		Name:         pgtype.Text{String: request.Name, Valid: request.Name != ""},
		Gender:       repo.NullGenderType{GenderType: request.Gender, Valid: true},
		Photo:        pgtype.Text{String: request.Photo, Valid: request.Photo != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		UserID:       pgtype.Int4{Int32: request.UserID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &model.Employee{
		ID:           result.ID,
		Name:         request.Name,
		NIP:          request.NIP,
		Gender:       request.Gender,
		Photo:        request.Photo,
		OccupationID: request.OccupationID,
		UserID:       request.UserID,
	}, nil
}

func (s *EmployeeUseCase) Delete(ctx context.Context, employeeId int32) (*model.Employee, error) {
	result, err := s.store.DeleteEmployee(ctx, employeeId)
	if err != nil {
		return nil, err
	}

	return &model.Employee{
		ID:           result.ID,
		Name:         result.Name,
		NIP:          result.Nip.String,
		Gender:       result.Gender,
		Photo:        result.Photo.String,
		OccupationID: result.OccupationID,
		UserID:       result.UserID.Int32,
	}, nil
}
