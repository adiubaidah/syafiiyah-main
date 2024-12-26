package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type SantriUseCase interface {
	CreateSantri(ctx context.Context, request *model.CreateSantriRequest) (*model.SantriResponse, error)
	ListSantri(ctx context.Context, request *model.ListSantriRequest) (*[]model.SantriCompleteResponse, error)
	CountSantri(ctx context.Context, request *model.ListSantriRequest) (int64, error)
	GetSantri(ctx context.Context, santriId int32) (*model.SantriCompleteResponse, error)
	UpdateSantri(ctx context.Context, request *model.UpdateSantriRequest, santriId int32) (*model.SantriResponse, error)
	DeleteSantri(ctx context.Context, santriId int32) (*model.SantriResponse, error)
}

type santriService struct {
	store repo.Store
}

func NewSantriUseCase(store repo.Store) SantriUseCase {
	return &santriService{store: store}
}

func (c *santriService) CreateSantri(ctx context.Context, request *model.CreateSantriRequest) (*model.SantriResponse, error) {
	isActive, err := strconv.ParseBool(request.IsActive)
	if err != nil {
		return nil, err
	}
	createdSantri, err := c.store.CreateSantri(ctx, repo.CreateSantriParams{
		Nis:          pgtype.Text{String: request.Nis, Valid: true},
		Name:         request.Name,
		IsActive:     pgtype.Bool{Bool: isActive, Valid: true},
		Generation:   request.Generation,
		Photo:        pgtype.Text{String: request.Photo, Valid: request.Photo != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		ParentID:     pgtype.Int4{Int32: request.ParentID, Valid: request.ParentID != 0},
		Gender:       request.Gender,
	})
	if err != nil {
		return nil, err
	}

	return &model.SantriResponse{
		ID: createdSantri.ID,

		Nis:          createdSantri.Nis.String,
		Name:         createdSantri.Name,
		Gender:       createdSantri.Gender,
		IsActive:     createdSantri.IsActive.Bool,
		Generation:   createdSantri.Generation,
		Photo:        createdSantri.Photo.String,
		OccupationID: createdSantri.OccupationID.Int32,
		ParentID:     createdSantri.ParentID.Int32,
	}, nil
}

func (c *santriService) ListSantri(ctx context.Context, request *model.ListSantriRequest) (*[]model.SantriCompleteResponse, error) {
	var result []model.SantriCompleteResponse
	offset := (request.Page - 1) * request.Limit
	if offset < 0 {
		offset = 0
	}

	arg := repo.ListSantriParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		Generation:   pgtype.Int4{Int32: request.Generation, Valid: request.Generation != 0},
		OffsetNumber: offset,
		LimitNumber:  request.Limit,
		IsActive:     pgtype.Bool{Bool: request.IsActive == 1, Valid: request.IsActive != 0},
		OrderBy:      repo.NullSantriOrderBy{SantriOrderBy: repo.SantriOrderBy(request.Order), Valid: request.Order != ""},
	}
	santris, err := c.store.ListSantri(ctx, arg)
	if err != nil {
		return nil, err
	}

	for _, santri := range santris {
		result = append(result, model.SantriCompleteResponse{
			ID:           santri.ID,
			Nis:          santri.Nis.String,
			Name:         santri.Name,
			Gender:       santri.Gender,
			IsActive:     santri.IsActive.Bool,
			Generation:   santri.Generation,
			Photo:        santri.Photo.String,
			OccupationID: santri.OccupationID.Int32,
			ParentID:     santri.ParentID.Int32,

			Occupation: model.SantriOccupation{
				ID:   santri.OccupationID.Int32,
				Name: santri.OccupationName.String,
			},
			Parent: model.SantriParent{
				ID:   santri.ParentID.Int32,
				Name: santri.ParentName.String,
			},
		})
	}

	return &result, nil
}

func (c *santriService) CountSantri(ctx context.Context, request *model.ListSantriRequest) (int64, error) {

	arg := repo.CountSantriParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		Generation:   pgtype.Int4{Int32: request.Generation, Valid: request.Generation != 0},
		IsActive:     pgtype.Bool{Bool: request.IsActive == 1, Valid: request.IsActive != 0},
	}

	count, err := c.store.CountSantri(ctx, arg)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *santriService) GetSantri(ctx context.Context, santriId int32) (*model.SantriCompleteResponse, error) {
	santri, err := c.store.GetSantri(ctx, santriId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri not found")
		}
		return nil, err
	}
	return &model.SantriCompleteResponse{
		ID:           santri.ID,
		Nis:          santri.Nis.String,
		Name:         santri.Name,
		Gender:       santri.Gender,
		IsActive:     santri.IsActive.Bool,
		Generation:   santri.Generation,
		Photo:        santri.Photo.String,
		OccupationID: santri.OccupationID.Int32,
		ParentID:     santri.ParentID.Int32,
		Occupation: model.SantriOccupation{
			ID:   santri.OccupationID.Int32,
			Name: santri.OccupationName.String,
		},
		Parent: model.SantriParent{
			ID:   santri.ParentID.Int32,
			Name: santri.ParentName.String,
		},
	}, nil
}

func (c *santriService) UpdateSantri(ctx context.Context, request *model.UpdateSantriRequest, santriId int32) (*model.SantriResponse, error) {
	isActive, err := strconv.ParseBool(request.IsActive)
	if err != nil {
		return nil, err
	}
	createdSantri, err := c.store.UpdateSantri(ctx, repo.UpdateSantriParams{
		ID:           santriId,
		Nis:          pgtype.Text{String: request.Nis, Valid: true},
		Name:         pgtype.Text{String: request.Name, Valid: request.Name != ""},
		IsActive:     pgtype.Bool{Bool: isActive, Valid: true},
		Generation:   pgtype.Int4{Int32: request.Generation, Valid: request.Generation != 0},
		Photo:        pgtype.Text{String: request.Photo, Valid: request.Photo != ""},
		OccupationID: pgtype.Int4{Int32: request.OccupationID, Valid: request.OccupationID != 0},
		ParentID:     pgtype.Int4{Int32: request.ParentID, Valid: request.ParentID != 0},
		Gender:       repo.NullGenderType{GenderType: request.Gender, Valid: true},
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri not found")
		}
		return nil, err
	}

	return &model.SantriResponse{
		ID:           createdSantri.ID,
		Nis:          createdSantri.Nis.String,
		Name:         createdSantri.Name,
		Gender:       createdSantri.Gender,
		IsActive:     createdSantri.IsActive.Bool,
		Generation:   createdSantri.Generation,
		Photo:        createdSantri.Photo.String,
		OccupationID: createdSantri.OccupationID.Int32,
		ParentID:     createdSantri.ParentID.Int32,
	}, nil
}

func (c *santriService) DeleteSantri(ctx context.Context, santriId int32) (*model.SantriResponse, error) {
	santri, err := c.store.DeleteSantri(ctx, santriId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Santri not found")
		}
		return nil, err
	}
	return &model.SantriResponse{
		ID:           santri.ID,
		Nis:          santri.Nis.String,
		Name:         santri.Name,
		Gender:       santri.Gender,
		IsActive:     santri.IsActive.Bool,
		Generation:   santri.Generation,
		Photo:        santri.Photo.String,
		OccupationID: santri.OccupationID.Int32,
		ParentID:     santri.ParentID.Int32,
	}, nil

}
