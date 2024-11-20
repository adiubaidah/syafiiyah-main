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

type HolidayUseCase interface {
	CreateHoliday(ctx context.Context, request *model.CreateHolidayRequest) (model.HolidayResponse, error)
	ListHolidays(ctx context.Context, request *model.ListHolidayRequest) ([]model.HolidayResponse, error)
	UpdateHoliday(ctx context.Context, request *model.UpdateHolidayRequest, holidayId int32) (model.HolidayResponse, error)
	DeleteHoliday(ctx context.Context, holidayId int32) (model.HolidayResponse, error)
}

type holidayService struct {
	store db.Store
}

func NewHolidayUseCase(store db.Store) HolidayUseCase {
	return &holidayService{store: store}
}

func (s *holidayService) CreateHoliday(ctx context.Context, request *model.CreateHolidayRequest) (model.HolidayResponse, error) {
	sqlStore := s.store.(*db.SQLStore)

	argCreateDates := make([]db.CreateHolidayDatesParams, 0)
	for _, date := range request.Dates {
		parsedDate, err := util.ParseDate(date)
		if err != nil {
			return model.HolidayResponse{}, exception.NewValidationError(err.Error())
		}
		argCreateDates = append(argCreateDates, db.CreateHolidayDatesParams{
			Date: pgtype.Date{Time: parsedDate, Valid: true},
		})
	}

	holiday, err := sqlStore.CreateHolidayWithDates(ctx, db.CreateHolidayParams{
		Name:        request.Name,
		Color:       pgtype.Text{String: request.Color, Valid: request.Color != ""},
		Description: pgtype.Text{String: request.Description, Valid: request.Description != ""},
	}, argCreateDates)
	if err != nil {

		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return model.HolidayResponse{}, exception.NewUniqueViolationError("date", err)
		}
		return model.HolidayResponse{}, err
	}

	return model.HolidayResponse{
		ID:    holiday.ID,
		Name:  holiday.Name,
		Color: holiday.Color.String,
		Dates: request.Dates,
	}, err
}

func (s *holidayService) ListHolidays(ctx context.Context, request *model.ListHolidayRequest) ([]model.HolidayResponse, error) {
	holidays, err := s.store.ListHolidays(ctx, db.ListHolidaysParams{
		Month: pgtype.Int4{Int32: request.Month, Valid: request.Month != 0},
		Year:  pgtype.Int4{Int32: request.Year, Valid: request.Year != 0},
	})
	if err != nil {
		return nil, err
	}

	holidayMap := make(map[int32]*model.HolidayResponse)
	for _, holiday := range holidays {
		if _, exists := holidayMap[holiday.ID]; !exists {
			holidayMap[holiday.ID] = &model.HolidayResponse{
				ID:          holiday.ID,
				Name:        holiday.Name,
				Color:       holiday.Color.String,
				Description: holiday.Description.String,
			}
		}
		holidayMap[holiday.ID].Dates = append(holidayMap[holiday.ID].Dates, holiday.HolidayDate.Time.Format("2006-01-02"))

	}
	var holidayResponses []model.HolidayResponse
	for _, holiday := range holidayMap {
		holidayResponses = append(holidayResponses, *holiday)
	}

	return holidayResponses, nil
}

func (s *holidayService) UpdateHoliday(ctx context.Context, request *model.UpdateHolidayRequest, holidayId int32) (model.HolidayResponse, error) {
	sqlStore := s.store.(*db.SQLStore)
	argCreateDates := make([]db.CreateHolidayDatesParams, 0)
	for _, date := range request.Dates {
		parsedDate, err := util.ParseDate(date)
		if err != nil {
			return model.HolidayResponse{}, exception.NewValidationError(err.Error())
		}
		argCreateDates = append(argCreateDates, db.CreateHolidayDatesParams{
			Date: pgtype.Date{Time: parsedDate, Valid: true},
		})
	}
	holiday, err := sqlStore.UpdateHolidayWithDates(ctx, holidayId, db.UpdateHolidayParams{
		ID:          holidayId,
		Name:        pgtype.Text{String: request.Name, Valid: request.Name != ""},
		Color:       pgtype.Text{String: request.Color, Valid: request.Color != ""},
		Description: pgtype.Text{String: request.Description, Valid: request.Description != ""},
	}, argCreateDates)
	if err != nil {
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return model.HolidayResponse{}, exception.NewUniqueViolationError("date", err)
		}

		if errors.Is(err, exception.ErrNotFound) {
			return model.HolidayResponse{}, exception.NewNotFoundError("Holiday not found")
		}

		return model.HolidayResponse{}, err
	}

	return model.HolidayResponse{
		ID:    holiday.ID,
		Name:  holiday.Name,
		Color: holiday.Color.String,
		Dates: request.Dates,
	}, err
}

func (s *holidayService) DeleteHoliday(ctx context.Context, holidayId int32) (model.HolidayResponse, error) {
	holiday, err := s.store.DeleteHoliday(ctx, holidayId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return model.HolidayResponse{}, exception.NewNotFoundError("Holiday not found")
		}
		return model.HolidayResponse{}, err
	}

	return model.HolidayResponse{
		ID:    holiday.ID,
		Name:  holiday.Name,
		Color: holiday.Color.String,
		Dates: nil,
	}, nil
}
