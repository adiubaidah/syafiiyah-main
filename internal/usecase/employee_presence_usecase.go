package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/exception"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmployeePresenceUseCase struct {
	store repo.Store
}

func NewEmployeePresenceUseCase(store repo.Store) *EmployeePresenceUseCase {
	return &EmployeePresenceUseCase{
		store: store,
	}
}

func (s *EmployeePresenceUseCase) CreatePresence(ctx context.Context, request *model.CreateEmployeePresenceRequest) (*model.EmployeePresenceResponse, error) {

	getEmployee, err := s.store.GetEmployeeByID(ctx, request.EmployeeID)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri not found")
		}
		return nil, err
	}
	result, err := s.store.CreateEmployeePresence(ctx, repo.CreateEmployeePresenceParams{
		ScheduleID:           request.ScheduleID,
		EmployeeID:           request.EmployeeID,
		ScheduleName:         request.ScheduleName,
		Type:                 request.Type,
		Notes:                pgtype.Text{String: request.Notes, Valid: request.Notes != ""},
		CreatedBy:            request.CreatedBy,
		EmployeePermissionID: pgtype.Int4{Int32: request.EmployeePermissionID, Valid: request.EmployeePermissionID != 0},
	})
	if err != nil {
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return nil, exception.NewUniqueViolationError("employee Presence today already exist", err)
		}
		return nil, err
	}

	return &model.EmployeePresenceResponse{
		ID:                   result.ID,
		Type:                 result.Type,
		EmployeeID:           result.EmployeeID,
		CreatedAt:            result.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		Notes:                result.Notes.String,
		EmployeePermissionID: result.EmployeeID,
		Schedule: model.IdAndName{
			Id:   result.ScheduleID,
			Name: result.ScheduleName,
		},
		Employee: model.IdAndName{
			Id:   result.EmployeeID,
			Name: getEmployee.Name,
		},
	}, nil

}

func (s *EmployeePresenceUseCase) BulkCreatePresence(ctx context.Context, args []repo.CreateEmployeePresencesParams) (int64, error) {

	if len(args) == 0 {
		return 0, exception.NewValidationError("Santri presence data is empty")
	}

	affected, err := s.store.CreateEmployeePresences(ctx, args)
	if err != nil {
		return 0, err
	}
	return affected, nil

}

func (s *EmployeePresenceUseCase) ListPresences(ctx context.Context, request *model.ListSantriPresenceRequest) (*[]model.EmployeePresenceResponse, error) {

	var fromDate, toDate time.Time
	var err error
	if request.From != "" {
		fromDate, err = util.ParseDate(request.From)
		if err != nil {
			return nil, exception.NewValidationError("From date is not valid")
		}
	}

	if request.To != "" {
		toDate, err = util.ParseDate(request.To)
		if err != nil {
			return nil, exception.NewValidationError("To date is not valid")
		}
	}

	results, err := s.store.ListEmployeePresences(ctx, repo.ListEmployeePresencesParams{
		EmployeeID:   pgtype.Int4{Int32: request.SantriID, Valid: request.SantriID != 0},
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		Type:         repo.NullPresenceType{PresenceType: request.Type, Valid: request.Type != ""},
		ScheduleID:   pgtype.Int4{Int32: request.ScheduleID, Valid: request.ScheduleID != 0},
		FromDate:     pgtype.Date{Time: fromDate, Valid: request.From != ""},
		ToDate:       pgtype.Date{Time: toDate, Valid: request.To != ""},
		OffsetNumber: (request.Page - 1) * request.Limit,
		LimitNumber:  request.Limit,
	})
	if err != nil {
		return nil, err
	}

	var response []model.EmployeePresenceResponse
	for _, employeePresence := range results {
		response = append(response, model.EmployeePresenceResponse{
			ID:                   employeePresence.ID,
			Type:                 employeePresence.Type,
			EmployeeID:           employeePresence.EmployeeID,
			CreatedAt:            employeePresence.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			Notes:                employeePresence.Notes.String,
			EmployeePermissionID: employeePresence.EmployeePermissionID.Int32,
			Schedule: model.IdAndName{
				Id:   employeePresence.ScheduleID,
				Name: employeePresence.ScheduleName,
			},
			Employee: model.IdAndName{
				Id:   employeePresence.EmployeeID,
				Name: employeePresence.EmployeeName,
			},
		})
	}

	return &response, nil
}

func (s *EmployeePresenceUseCase) CountPresences(ctx context.Context, request *model.ListEmployeePresenceRequest) (int64, error) {
	var fromDate, toDate time.Time
	var err error
	if request.From != "" {
		fromDate, err = util.ParseDate(request.From)
		if err != nil {
			return 0, exception.NewValidationError("From date is not valid")
		}
	}

	if request.To != "" {
		toDate, err = util.ParseDate(request.To)
		if err != nil {
			return 0, exception.NewValidationError("To date is not valid")
		}
	}

	count, err := s.store.CountEmployeePresences(ctx, repo.CountEmployeePresencesParams{
		EmployeeID: pgtype.Int4{Int32: request.EmployeeID, Valid: request.EmployeeID != 0},
		Q:          pgtype.Text{String: request.Q, Valid: request.Q != ""},
		Type:       repo.NullPresenceType{PresenceType: request.Type, Valid: request.Type != ""},
		ScheduleID: pgtype.Int4{Int32: request.ScheduleID, Valid: request.ScheduleID != 0},
		FromDate:   pgtype.Date{Time: fromDate, Valid: request.From != ""},
		ToDate:     pgtype.Date{Time: toDate, Valid: request.To != ""},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *EmployeePresenceUseCase) ListMissingPresences(ctx context.Context, request *model.ListMissingEmployeePresenceRequest) (*[]model.IdAndName, error) {
	missingSantriPresences, err := s.store.ListMissingEmployeePresences(ctx, repo.ListMissingEmployeePresencesParams{
		Date:       pgtype.Date{Time: request.Time, Valid: true},
		ScheduleID: pgtype.Int4{Int32: request.ScheduleID, Valid: request.ScheduleID != 0},
	})
	if err != nil {
		return nil, err
	}

	var response []model.IdAndName
	for _, missingSantriPresence := range missingSantriPresences {
		response = append(response, model.IdAndName{
			Id:   missingSantriPresence.ID,
			Name: missingSantriPresence.Name,
		})
	}

	return &response, nil
}

func (s *EmployeePresenceUseCase) UpdatePresence(ctx context.Context, request *model.UpdateEmployeePresenceRequest, santriPresenceID int32) (*model.EmployeePresenceResponse, error) {

	getEmployee, err := s.store.GetEmployeeByID(ctx, request.EmployeeID)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri not found")
		}
		return nil, err
	}

	result, err := s.store.UpdateEmployeePresence(ctx, repo.UpdateEmployeePresenceParams{
		ID:                   santriPresenceID,
		ScheduleID:           pgtype.Int4{Int32: request.ScheduleID, Valid: request.ScheduleID != 0},
		ScheduleName:         pgtype.Text{String: request.ScheduleName, Valid: request.ScheduleName != ""},
		Type:                 repo.NullPresenceType{PresenceType: request.Type, Valid: request.Type != ""},
		Notes:                pgtype.Text{String: request.Notes, Valid: request.Notes != ""},
		EmployeePermissionID: pgtype.Int4{Int32: request.EmployeePermissionID, Valid: request.EmployeePermissionID != 0},
		EmployeeID:           pgtype.Int4{Int32: request.EmployeeID, Valid: request.EmployeeID != 0},
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Presence not found")
		}
		return nil, err
	}

	return &model.EmployeePresenceResponse{
		ID:                   result.ID,
		Type:                 result.Type,
		EmployeeID:           result.EmployeeID,
		CreatedAt:            result.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		Notes:                result.Notes.String,
		EmployeePermissionID: result.EmployeePermissionID.Int32,
		Schedule: model.IdAndName{
			Id:   result.ScheduleID,
			Name: result.ScheduleName,
		},
		Employee: model.IdAndName{
			Id:   result.EmployeeID,
			Name: getEmployee.Name,
		},
	}, nil
}

func (s *EmployeePresenceUseCase) DeletePresence(ctx context.Context, santriPresenceID int32) (*model.EmployeePresenceResponse, error) {
	deleted, err := s.store.DeleteEmployeePresence(ctx, santriPresenceID)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Presence not found")
		}
		return nil, err
	}

	return &model.EmployeePresenceResponse{
		ID:                   deleted.ID,
		Type:                 deleted.Type,
		EmployeeID:           deleted.EmployeeID,
		CreatedAt:            deleted.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		Notes:                deleted.Notes.String,
		EmployeePermissionID: deleted.EmployeePermissionID.Int32,
		Schedule: model.IdAndName{
			Id:   deleted.ScheduleID,
			Name: deleted.ScheduleName,
		},
		Employee: model.IdAndName{
			Id: deleted.EmployeeID,
		},
	}, nil
}
