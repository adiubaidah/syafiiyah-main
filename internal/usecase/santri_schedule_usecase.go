package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type SantriScheduleUseCase interface {
	CreateSantriSchedule(ctx context.Context, request *model.CreateSantriScheduleRequest) (*model.SantriScheduleResponse, error)
	ListSantriSchedule(ctx context.Context) (*[]model.SantriScheduleResponse, error)
	GetSantriSchedule(ctx context.Context, time time.Time) (*model.SantriScheduleResponse, error)
	UpdateSantriSchedule(ctx context.Context, request *model.UpdateSantriScheduleRequest, santriScheduleId int32) (*model.SantriScheduleResponse, error)
	DeleteSantriSchedule(ctx context.Context, santriScheduleId int32) (*model.SantriScheduleResponse, error)
}

type santriScheduleService struct {
	store db.Store
}

func NewSantriScheduleUseCase(store db.Store) SantriScheduleUseCase {
	return &santriScheduleService{store: store}
}

func (c *santriScheduleService) CreateSantriSchedule(ctx context.Context, request *model.CreateSantriScheduleRequest) (*model.SantriScheduleResponse, error) {
	startPresence, err := util.ParseTime(request.StartPresence)
	if err != nil {
		return nil, err
	}
	startTime, err := util.ParseTime(request.StartTime)
	if err != nil {
		return nil, err
	}
	finishTime, err := util.ParseTime(request.FinishTime)
	if err != nil {
		return nil, err
	}

	crashStartPresence, err := c.store.ListSantriSchedules(ctx, util.ConvertToPgxTime(startPresence))
	if err != nil {
		return nil, err
	}

	if len(crashStartPresence) > 0 {
		var crashedName []string
		for _, presence := range crashStartPresence {
			crashedName = append(crashedName, presence.Name)
		}

		return nil, exception.NewValidationError(fmt.Sprintf("start presence crash with %s", strings.Join(crashedName, ", ")))
	}

	crashFinishTime, err := c.store.ListSantriSchedules(ctx, util.ConvertToPgxTime(finishTime))
	if err != nil {
		return nil, err
	}

	if len(crashFinishTime) > 0 {
		var crashedName []string
		for _, presence := range crashFinishTime {
			crashedName = append(crashedName, presence.Name)
		}

		return nil, exception.NewValidationError(fmt.Sprintf("finish time crash with %s", strings.Join(crashedName, ", ")))
	}

	if !startTime.After(startPresence) {
		return nil, exception.NewValidationError("start_time must be after start_presence")
	}
	if !finishTime.After(startTime) {
		return nil, exception.NewValidationError("finish_time must be after start_presence")
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
			return nil, exception.NewUniqueViolationError("start presence, start time or finish time has been crashed", err)
		}

		return nil, exception.NewDatabaseError("CreateSantriSchedule", err)
	}

	return &model.SantriScheduleResponse{
		ID:            createdSantriSchedule.ID,
		Name:          createdSantriSchedule.Name,
		Description:   createdSantriSchedule.Description.String,
		StartPresence: util.ConvertToTime(createdSantriSchedule.StartPresence),
		StartTime:     util.ConvertToTime(createdSantriSchedule.StartTime),
		FinishTime:    util.ConvertToTime(createdSantriSchedule.FinishTime),
	}, nil
}

func (c *santriScheduleService) ListSantriSchedule(ctx context.Context) (*[]model.SantriScheduleResponse, error) {
	santriSchedules, err := c.store.ListSantriSchedules(ctx, pgtype.Time{Valid: false})
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

	return &response, nil
}

func (c *santriScheduleService) GetSantriSchedule(ctx context.Context, time time.Time) (*model.SantriScheduleResponse, error) {
	schedule, err := c.store.GetSantriSchedule(ctx, util.ConvertToPgxTime(time))
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Schedule not found")
		}
		return nil, err
	}
	return &model.SantriScheduleResponse{
		ID:            schedule.ID,
		Name:          schedule.Name,
		Description:   schedule.Description.String,
		StartPresence: util.ConvertToTime(schedule.StartPresence),
		StartTime:     util.ConvertToTime(schedule.StartTime),
		FinishTime:    util.ConvertToTime(schedule.FinishTime),
	}, nil
}

func (c *santriScheduleService) UpdateSantriSchedule(ctx context.Context, request *model.UpdateSantriScheduleRequest, santriScheduleId int32) (*model.SantriScheduleResponse, error) {
	startPresence, err := util.ParseTime(request.StartPresence)
	if err != nil {
		return nil, err
	}
	startTime, err := util.ParseTime(request.StartTime)
	if err != nil {
		return nil, err
	}
	finishTime, err := util.ParseTime(request.FinishTime)
	if err != nil {
		return nil, err
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
			return nil, exception.NewNotFoundError("Santri Schedule not found")
		}
		return nil, err
	}

	return &model.SantriScheduleResponse{
		ID:            updatedSantriSchedule.ID,
		Name:          updatedSantriSchedule.Name,
		Description:   updatedSantriSchedule.Description.String,
		StartPresence: util.ConvertToTime(updatedSantriSchedule.StartPresence),
		StartTime:     util.ConvertToTime(updatedSantriSchedule.StartTime),
		FinishTime:    util.ConvertToTime(updatedSantriSchedule.FinishTime),
	}, nil
}

func (c *santriScheduleService) DeleteSantriSchedule(ctx context.Context, santriScheduleId int32) (*model.SantriScheduleResponse, error) {
	deletedSantriSchedule, err := c.store.DeleteSantriSchedule(ctx, santriScheduleId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri Schedule not found")
		}
		return nil, err
	}

	return &model.SantriScheduleResponse{
		ID:            deletedSantriSchedule.ID,
		Name:          deletedSantriSchedule.Name,
		Description:   deletedSantriSchedule.Description.String,
		StartPresence: util.ConvertToTime(deletedSantriSchedule.StartPresence),
		StartTime:     util.ConvertToTime(deletedSantriSchedule.StartTime),
		FinishTime:    util.ConvertToTime(deletedSantriSchedule.FinishTime),
	}, nil
}
