package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type SantriPresenceUseCase interface {
	CreateSantriPresence(ctx context.Context, request *model.CreateSantriPresenceRequest) (*model.SantriPresenceResponse, error)
	BulkCreateSantriPresence(ctx context.Context, args []repo.CreateSantriPresencesParams) (int64, error)
	ListSantriPresences(ctx context.Context, request *model.ListSantriPresenceRequest) (*[]model.SantriPresenceResponse, error)
	ListMissingSantriPresences(ctx context.Context, request *model.ListMissingSantriPresenceRequest) (*[]model.IdAndName, error)
	CountSantriPresences(ctx context.Context, request *model.ListSantriPresenceRequest) (int64, error)
	UpdateSantriPresence(ctx context.Context, request *model.UpdateSantriPresenceRequest, santriPresenceID int32) (*model.SantriPresenceResponse, error)
	DeleteSantriPresence(ctx context.Context, santriPresenceID int32) (*model.SantriPresenceResponse, error)
}

type santriPresenceService struct {
	store repo.Store
}

func NewSantriPresenceUseCase(store repo.Store) SantriPresenceUseCase {
	return &santriPresenceService{
		store: store,
	}
}

func (s *santriPresenceService) CreateSantriPresence(ctx context.Context, request *model.CreateSantriPresenceRequest) (*model.SantriPresenceResponse, error) {

	getSantri, err := s.store.GetSantri(ctx, request.SantriID)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri not found")
		}
		return nil, err
	}
	createdSantriPresence, err := s.store.CreateSantriPresence(ctx, repo.CreateSantriPresenceParams{
		ScheduleID:         request.ScheduleID,
		SantriID:           request.SantriID,
		ScheduleName:       request.ScheduleName,
		Type:               request.Type,
		Notes:              pgtype.Text{String: request.Notes, Valid: request.Notes != ""},
		CreatedBy:          request.CreatedBy,
		SantriPermissionID: pgtype.Int4{Int32: request.SantriPermissionID, Valid: request.SantriPermissionID != 0},
	})
	if err != nil {
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return nil, exception.NewUniqueViolationError("Santri Presence today already exist", err)
		}
		return nil, err
	}

	return &model.SantriPresenceResponse{
		ID:                 createdSantriPresence.ID,
		Type:               createdSantriPresence.Type,
		SantriID:           createdSantriPresence.SantriID,
		CreatedAt:          createdSantriPresence.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		Notes:              createdSantriPresence.Notes.String,
		SantriPermissionID: createdSantriPresence.SantriPermissionID.Int32,
		Schedule: model.IdAndName{
			Id:   createdSantriPresence.ScheduleID,
			Name: createdSantriPresence.ScheduleName,
		},
		Santri: model.IdAndName{
			Id:   createdSantriPresence.SantriID,
			Name: getSantri.Name,
		},
	}, nil

}

func (s *santriPresenceService) BulkCreateSantriPresence(ctx context.Context, args []repo.CreateSantriPresencesParams) (int64, error) {

	if len(args) == 0 {
		return 0, exception.NewValidationError("Santri presence data is empty")
	}

	affected, err := s.store.CreateSantriPresences(ctx, args)
	if err != nil {
		return 0, err
	}
	return affected, nil

}

func (s *santriPresenceService) ListSantriPresences(ctx context.Context, request *model.ListSantriPresenceRequest) (*[]model.SantriPresenceResponse, error) {

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

	santriPresences, err := s.store.ListSantriPresences(ctx, repo.ListSantriPresencesParams{
		SantriID:     pgtype.Int4{Int32: request.SantriID, Valid: request.SantriID != 0},
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

	var response []model.SantriPresenceResponse
	for _, santriPresence := range santriPresences {
		response = append(response, model.SantriPresenceResponse{
			ID:                 santriPresence.ID,
			Type:               santriPresence.Type,
			SantriID:           santriPresence.SantriID,
			CreatedAt:          santriPresence.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			Notes:              santriPresence.Notes.String,
			SantriPermissionID: santriPresence.SantriPermissionID.Int32,
			Schedule: model.IdAndName{
				Id:   santriPresence.ScheduleID,
				Name: santriPresence.ScheduleName,
			},
			Santri: model.IdAndName{
				Id:   santriPresence.SantriID,
				Name: santriPresence.SantriName,
			},
		})
	}

	return &response, nil
}

func (s *santriPresenceService) CountSantriPresences(ctx context.Context, request *model.ListSantriPresenceRequest) (int64, error) {
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

	count, err := s.store.CountSantriPresences(ctx, repo.CountSantriPresencesParams{
		SantriID:   pgtype.Int4{Int32: request.SantriID, Valid: request.SantriID != 0},
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

func (s *santriPresenceService) ListMissingSantriPresences(ctx context.Context, request *model.ListMissingSantriPresenceRequest) (*[]model.IdAndName, error) {
	missingSantriPresences, err := s.store.ListMissingSantriPresences(ctx, repo.ListMissingSantriPresencesParams{
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

func (s *santriPresenceService) UpdateSantriPresence(ctx context.Context, request *model.UpdateSantriPresenceRequest, santriPresenceID int32) (*model.SantriPresenceResponse, error) {

	getSantri, err := s.store.GetSantri(ctx, request.SantriID)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri not found")
		}
		return nil, err
	}

	updatedSantriPresence, err := s.store.UpdateSantriPresence(ctx, repo.UpdateSantriPresenceParams{
		ID:                 santriPresenceID,
		ScheduleID:         pgtype.Int4{Int32: request.ScheduleID, Valid: request.ScheduleID != 0},
		Type:               repo.NullPresenceType{PresenceType: request.Type, Valid: request.Type != ""},
		Notes:              pgtype.Text{String: request.Notes, Valid: request.Notes != ""},
		SantriPermissionID: pgtype.Int4{Int32: request.SantriPermissionID, Valid: request.SantriPermissionID != 0},
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Presence not found")
		}
		return nil, err
	}

	return &model.SantriPresenceResponse{
		ID:                 updatedSantriPresence.ID,
		Type:               updatedSantriPresence.Type,
		SantriID:           updatedSantriPresence.SantriID,
		CreatedAt:          updatedSantriPresence.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		Notes:              updatedSantriPresence.Notes.String,
		SantriPermissionID: updatedSantriPresence.SantriPermissionID.Int32,
		Schedule: model.IdAndName{
			Id:   updatedSantriPresence.ScheduleID,
			Name: updatedSantriPresence.ScheduleName,
		},
		Santri: model.IdAndName{
			Id:   updatedSantriPresence.SantriID,
			Name: getSantri.Name,
		},
	}, nil
}

func (s *santriPresenceService) DeleteSantriPresence(ctx context.Context, santriPresenceID int32) (*model.SantriPresenceResponse, error) {
	deleted, err := s.store.DeleteSantriPresence(ctx, santriPresenceID)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Presence not found")
		}
		return nil, err
	}

	return &model.SantriPresenceResponse{
		ID:                 deleted.ID,
		Type:               deleted.Type,
		SantriID:           deleted.SantriID,
		CreatedAt:          deleted.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		Notes:              deleted.Notes.String,
		SantriPermissionID: deleted.SantriPermissionID.Int32,
		Schedule: model.IdAndName{
			Id:   deleted.ScheduleID,
			Name: deleted.ScheduleName,
		},
		Santri: model.IdAndName{
			Id: deleted.SantriID,
		},
	}, nil
}
