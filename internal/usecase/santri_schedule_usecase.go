package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type SantriScheduleUseCase interface {
	CreateSantriSchedule(ctx context.Context, request *model.CreateSantriScheduleRequest) (model.SantriScheduleResponse, error)
	ListSantriSchedule(ctx context.Context) ([]model.SantriScheduleResponse, error)
	UpdateSantriSchedule(ctx context.Context, request *model.UpdateSantriScheduleRequest, santriScheduleId int32) (model.SantriScheduleResponse, error)
	DeleteSantriSchedule(ctx context.Context, santriScheduleId int32) (model.SantriScheduleResponse, error)
}

type santriScheduleService struct {
	store db.Store
}

func NewSantriScheduleUseCase(store db.Store) SantriScheduleUseCase {
	return &santriScheduleService{store: store}
}

func (c *santriScheduleService) CreateSantriSchedule(ctx context.Context, request *model.CreateSantriScheduleRequest) (model.SantriScheduleResponse, error) {
	startPresence, err := util.ParseTime(request.StartPresence)
	if err != nil {
		return model.SantriScheduleResponse{}, err
	}
	startTime, err := util.ParseTime(request.StartTime)
	if err != nil {
		return model.SantriScheduleResponse{}, err
	}
	finishTime, err := util.ParseTime(request.FinishTime)
	if err != nil {
		return model.SantriScheduleResponse{}, err
	}

	if !startTime.After(startPresence) {
		return model.SantriScheduleResponse{}, exception.NewValidationError("start_time must be after start_presence")
	}
	if !finishTime.After(startTime) {
		return model.SantriScheduleResponse{}, exception.NewValidationError("finish_time must be after start_presence")
	}

	createdSantriSchedule, err := c.store.CreateSantriSchedule(ctx, db.CreateSantriScheduleParams{
		Name:          request.Name,
		Description:   pgtype.Text{String: request.Description, Valid: request.Description != ""},
		StartPresence: util.ConvertToPgxTime(startPresence),
		StartTime:     util.ConvertToPgxTime(startTime),
		FinishTime:    util.ConvertToPgxTime(finishTime),
	})
	if err != nil {
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return model.SantriScheduleResponse{}, exception.NewUniqueViolationError("start presence, start time or finish time has been crashed", err)
		}

		return model.SantriScheduleResponse{}, exception.NewDatabaseError("CreateSantriSchedule", err)
	}

	return model.SantriScheduleResponse{
		ID:            createdSantriSchedule.ID,
		Name:          createdSantriSchedule.Name,
		Description:   createdSantriSchedule.Description.String,
		StartPresence: util.ConvertToTime(createdSantriSchedule.StartPresence),
		StartTime:     util.ConvertToTime(createdSantriSchedule.StartTime),
		FinishTime:    util.ConvertToTime(createdSantriSchedule.FinishTime),
	}, nil
}

func (c *santriScheduleService) ListSantriSchedule(ctx context.Context) ([]model.SantriScheduleResponse, error) {
	santriSchedules, err := c.store.ListSantriSchedules(ctx)
	if err != nil {
		return nil, err
	}

	var response []model.SantriScheduleResponse
	for _, santriSchedule := range santriSchedules {
		response = append(response, model.SantriScheduleResponse{
			ID:            santriSchedule.ID,
			Name:          santriSchedule.Name,
			Description:   santriSchedule.Description.String,
			StartPresence: util.ConvertToTime(santriSchedule.StartPresence),
			StartTime:     util.ConvertToTime(santriSchedule.StartTime),
			FinishTime:    util.ConvertToTime(santriSchedule.FinishTime),
		})
	}

	return response, nil
}

func (c *santriScheduleService) UpdateSantriSchedule(ctx context.Context, request *model.UpdateSantriScheduleRequest, santriScheduleId int32) (model.SantriScheduleResponse, error) {
	startPresence, err := util.ParseTime(request.StartPresence)
	if err != nil {
		return model.SantriScheduleResponse{}, err
	}
	startTime, err := util.ParseTime(request.StartTime)
	if err != nil {
		return model.SantriScheduleResponse{}, err
	}
	finishTime, err := util.ParseTime(request.FinishTime)
	if err != nil {
		return model.SantriScheduleResponse{}, err
	}

	updatedSantriSchedule, err := c.store.UpdateSantriSchedule(ctx, db.UpdateSantriScheduleParams{
		ID:            santriScheduleId,
		Name:          pgtype.Text{String: request.Name, Valid: request.Name != ""},
		Description:   pgtype.Text{String: request.Description, Valid: request.Description != ""},
		StartPresence: util.ConvertToPgxTime(startPresence),
		StartTime:     util.ConvertToPgxTime(startTime),
		FinishTime:    util.ConvertToPgxTime(finishTime),
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return model.SantriScheduleResponse{}, exception.NewNotFoundError("Santri Schedule not found")
		}
		return model.SantriScheduleResponse{}, err
	}

	return model.SantriScheduleResponse{
		ID:            updatedSantriSchedule.ID,
		Name:          updatedSantriSchedule.Name,
		Description:   updatedSantriSchedule.Description.String,
		StartPresence: util.ConvertToTime(updatedSantriSchedule.StartPresence),
		StartTime:     util.ConvertToTime(updatedSantriSchedule.StartTime),
		FinishTime:    util.ConvertToTime(updatedSantriSchedule.FinishTime),
	}, nil
}

func (c *santriScheduleService) DeleteSantriSchedule(ctx context.Context, santriScheduleId int32) (model.SantriScheduleResponse, error) {
	deletedSantriSchedule, err := c.store.DeleteSantriSchedule(ctx, santriScheduleId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return model.SantriScheduleResponse{}, exception.NewNotFoundError("Santri Schedule not found")
		}
		return model.SantriScheduleResponse{}, err
	}

	return model.SantriScheduleResponse{
		ID:            deletedSantriSchedule.ID,
		Name:          deletedSantriSchedule.Name,
		Description:   deletedSantriSchedule.Description.String,
		StartPresence: util.ConvertToTime(deletedSantriSchedule.StartPresence),
		StartTime:     util.ConvertToTime(deletedSantriSchedule.StartTime),
		FinishTime:    util.ConvertToTime(deletedSantriSchedule.FinishTime),
	}, nil
}
