package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type ParentUseCase interface {
	CreateParent(ctx context.Context, request *model.CreateParentRequest) (*model.ParentResponse, error)
	ListParents(ctx context.Context, request *model.ListParentRequest) (*[]model.ParentCompleteResponse, error)
	GetParent(ctx context.Context, parentId int32) (*model.ParentResponse, error)
	GetParentByUserID(ctx context.Context, userId int32) (*model.ParentResponse, error)
	CountParents(ctx context.Context, request *model.ListParentRequest) (int64, error)
	UpdateParent(ctx context.Context, request *model.UpdateParentRequest, parentId int32) (*model.ParentResponse, error)
	DeleteParent(ctx context.Context, parentId int32) (*model.ParentResponse, error)
}

type parentService struct {
	store repo.Store
}

func NewParentUseCase(store repo.Store) ParentUseCase {
	return &parentService{store: store}
}

func (c *parentService) CreateParent(ctx context.Context, request *model.CreateParentRequest) (*model.ParentResponse, error) {
	arg := repo.CreateParentParams{
		Name:           request.Name,
		Address:        request.Address,
		WhatsappNumber: pgtype.Text{String: request.WhatsappNumber, Valid: request.WhatsappNumber != ""},
		Gender:         request.Gender,
		Photo:          pgtype.Text{String: request.Photo, Valid: request.Photo != ""},
		UserID:         pgtype.Int4{Int32: request.UserID, Valid: request.UserID != 0},
	}

	createdParent, err := c.store.CreateParent(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &model.ParentResponse{
		ID:             createdParent.ID,
		Name:           createdParent.Name,
		Address:        createdParent.Address,
		WhatsappNumber: createdParent.WhatsappNumber.String,
		Gender:         string(createdParent.Gender),
		Photo:          createdParent.Photo.String,
		UserID:         createdParent.UserID.Int32,
	}, nil
}

func (c *parentService) ListParents(ctx context.Context, request *model.ListParentRequest) (*[]model.ParentCompleteResponse, error) {
	arg := repo.ListParentParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		HasUser:      pgtype.Bool{Bool: request.HasUser == 1, Valid: request.HasUser != 0},
		LimitNumber:  request.Limit,
		OffsetNumber: (request.Page - 1) * request.Limit,
		OrderBy:      repo.NullParentOrderBy{ParentOrderBy: repo.ParentOrderBy(request.Order), Valid: request.Order != ""},
	}

	parents, err := c.store.ListParents(ctx, arg)
	if err != nil {
		return nil, err
	}

	var result []model.ParentCompleteResponse
	for _, parent := range parents {
		result = append(result, model.ParentCompleteResponse{
			ParentResponse: model.ParentResponse{
				ID:             parent.ID,
				Name:           parent.Name,
				Address:        parent.Address,
				WhatsappNumber: parent.WhatsappNumber.String,
				Gender:         string(parent.Gender),
				Photo:          parent.Photo.String,
			},
			User: model.ParentUser{
				ID:       parent.UserID.Int32,
				Username: parent.Username.String,
			},
		})
	}
	return &result, nil
}

func (c *parentService) CountParents(ctx context.Context, request *model.ListParentRequest) (int64, error) {
	arg := repo.CountParentsParams{
		Q:       pgtype.Text{String: request.Q, Valid: request.Q != ""},
		HasUser: pgtype.Bool{Bool: request.HasUser == 1, Valid: request.HasUser != -1},
	}

	count, err := c.store.CountParents(ctx, arg)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *parentService) UpdateParent(ctx context.Context, request *model.UpdateParentRequest, parentId int32) (*model.ParentResponse, error) {
	arg := repo.UpdateParentParams{
		ID:             parentId,
		Name:           pgtype.Text{String: request.Name, Valid: request.Name != ""},
		Address:        pgtype.Text{String: request.Address, Valid: request.Address != ""},
		WhatsappNumber: pgtype.Text{String: request.WhatsappNumber, Valid: request.WhatsappNumber != ""},
		Gender:         repo.NullGenderType{GenderType: request.Gender, Valid: true},
		Photo:          pgtype.Text{String: request.Photo, Valid: request.Photo != ""},
		UserID:         pgtype.Int4{Int32: request.UserID, Valid: request.UserID != 0},
	}

	updatedParent, err := c.store.UpdateParent(ctx, arg)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Parent not found")
		}
		return nil, err
	}
	return &model.ParentResponse{
		ID:             updatedParent.ID,
		Name:           updatedParent.Name,
		Address:        updatedParent.Address,
		WhatsappNumber: updatedParent.WhatsappNumber.String,
		Gender:         string(updatedParent.Gender),
		Photo:          updatedParent.Photo.String,
		UserID:         updatedParent.UserID.Int32,
	}, nil
}

func (c *parentService) GetParent(ctx context.Context, parentId int32) (*model.ParentResponse, error) {
	parent, err := c.store.GetParent(ctx, parentId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Parent not found")
		}
		return nil, err
	}
	return &model.ParentResponse{
		ID:             parent.ID,
		Name:           parent.Name,
		Address:        parent.Address,
		WhatsappNumber: parent.WhatsappNumber.String,
		Gender:         string(parent.Gender),
		Photo:          parent.Photo.String,
		UserID:         parent.UserID.Int32,
	}, nil
}

func (c *parentService) GetParentByUserID(ctx context.Context, userID int32) (*model.ParentResponse, error) {
	parent, err := c.store.GetParentByUserId(ctx, pgtype.Int4{Int32: userID, Valid: true})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Parent not found")
		}
		return nil, err
	}
	return &model.ParentResponse{
		ID:             parent.ID,
		Name:           parent.Name,
		Address:        parent.Address,
		WhatsappNumber: parent.WhatsappNumber.String,
		Gender:         string(parent.Gender),
		Photo:          parent.Photo.String,
		UserID:         parent.UserID.Int32,
	}, nil
}

func (c *parentService) DeleteParent(ctx context.Context, parentId int32) (*model.ParentResponse, error) {
	parent, err := c.store.DeleteParent(ctx, parentId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Parent not found")
		}
		return nil, err
	}
	return &model.ParentResponse{
		ID:             parent.ID,
		Name:           parent.Name,
		Address:        parent.Address,
		WhatsappNumber: parent.WhatsappNumber.String,
		Gender:         string(parent.Gender),
		Photo:          parent.Photo.String,
		UserID:         parent.UserID.Int32,
	}, nil
}
