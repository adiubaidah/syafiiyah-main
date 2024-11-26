package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/jackc/pgx/v5/pgtype"
)

type SantriPresenceUseCase interface {
	CreateSantriPresence(ctx context.Context, request *model.CreateSantriPresenceRequest) (*model.SantriPresenceResponse, error)
}

type santriPresenceService struct {
	store db.Store
}

func NewSantriPresenceUseCase(store db.Store) SantriPresenceUseCase {
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
	createdSantriPresence, err := s.store.CreateSantriPresence(ctx, db.CreateSantriPresenceParams{
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
		ID:                 createdSantriPresence.ID.Int32,
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
