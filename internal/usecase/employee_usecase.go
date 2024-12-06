package usecase

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmployeeUseCase interface {
	CreateEmployee(ctx context.Context, request *model.CreateEmployeeRequest) (*model.Employee, error)
	ListEmployees(ctx context.Context, request *model.ListEmployeeRequest) (*[]model.EmployeeComplete, error)
	CountEmployees(ctx context.Context, request *model.ListEmployeeRequest) (int64, error)
	UpdateEmployee(ctx context.Context, request *model.UpdateEmployeeRequest, employeeId int32) (*model.Employee, error)
	DeleteEmployee(ctx context.Context, employeeId int32) (*model.Employee, error)
}

type employeeService struct {
	store db.Store
}

func NewEmployeeUseCase(store db.Store) EmployeeUseCase {
	return &employeeService{
		store: store,
	}
}

func (s *employeeService) CreateEmployee(ctx context.Context, request *model.CreateEmployeeRequest) (*model.Employee, error) {
	result, err := s.store.CreateEmployee(ctx, db.CreateEmployeeParams{
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

func (s *employeeService) ListEmployees(ctx context.Context, request *model.ListEmployeeRequest) (*[]model.EmployeeComplete, error) {
	employees, err := s.store.ListEmployees(ctx, db.ListEmployeesParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		HasUser:      pgtype.Bool{Bool: request.HasUser == 1, Valid: request.HasUser != 0},
		LimitNumber:  request.Limit,
		OffsetNumber: (request.Page - 1) * request.Limit,
		OrderBy:      db.NullEmployeeOrderBy{EmployeeOrderBy: request.Order, Valid: request.Order != ""},
	})
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
				Name: employee.Name,
			},
		})
	}

	return &result, nil

}

func (s *employeeService) CountEmployees(ctx context.Context, request *model.ListEmployeeRequest) (int64, error) {
	count, err := s.store.CountEmployees(ctx, db.CountEmployeesParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		HasUser:      pgtype.Bool{Bool: request.HasUser == 1, Valid: request.HasUser != 0},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, request *model.UpdateEmployeeRequest, employeeId int32) (*model.Employee, error) {
	result, err := s.store.UpdateEmployee(ctx, db.UpdateEmployeeParams{
		ID:           employeeId,
		Nip:          pgtype.Text{String: request.NIP, Valid: true}, // Nip can be null
		Name:         pgtype.Text{String: request.Name, Valid: request.Name != ""},
		Gender:       db.NullGenderType{GenderType: request.Gender, Valid: true},
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

func (s *employeeService) DeleteEmployee(ctx context.Context, employeeId int32) (*model.Employee, error) {
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
