package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type SmartCardUseCase interface {
	CreateSmartCard(ctx context.Context, request *model.SmartCardRequest) (*model.SmartCard, error)
	ListSmartCards(ctx context.Context, request *model.ListSmartCardRequest) (*[]model.SmartCardComplete, error)
	GetSmartCard(ctx context.Context, request *model.SmartCardRequest) (*model.SmartCardComplete, error)
	CountSmartCards(ctx context.Context, request *model.ListSmartCardRequest) (int64, error)
	UpdateSmartCard(ctx context.Context, request *model.UpdateSmartCardRequest, id int32) (*model.SmartCardComplete, error)
	DeleteSmartCard(ctx context.Context, id int32) (*model.SmartCard, error)
}

type service struct {
	store repo.Store
}

func NewSmartCardUseCase(store repo.Store) SmartCardUseCase {
	return &service{store: store}
}

func (c *service) CreateSmartCard(ctx context.Context, request *model.SmartCardRequest) (*model.SmartCard, error) {
	createdSmartCard, err := c.store.CreateSmartCard(ctx, repo.CreateSmartCardParams{
		Uid:        request.Uid,
		IsActive:   true,
		SantriID:   pgtype.Int4{Valid: false},
		EmployeeID: pgtype.Int4{Valid: false},
	})

	if err != nil {
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return nil, exception.NewUniqueViolationError("Smart Card already exists", err)
		}
		return nil, err
	}

	return &model.SmartCard{
		ID:        createdSmartCard.ID,
		Uid:       createdSmartCard.Uid,
		CreatedAt: createdSmartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		IsActive:  createdSmartCard.IsActive,
	}, nil
}

func (c *service) ListSmartCards(ctx context.Context, request *model.ListSmartCardRequest) (*[]model.SmartCardComplete, error) {
	listSmartCard, err := c.store.ListSmartCards(ctx, repo.ListSmartCardsParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
		CardOwner:    repo.NullCardOwner{CardOwner: request.CardOwner, Valid: request.CardOwner != ""},
		OffsetNumber: request.Limit * (request.Page - 1),
		LimitNumber:  request.Limit,
	})

	if err != nil {
		return nil, err
	}

	var result []model.SmartCardComplete

	for _, smartCard := range listSmartCard {

		ownerRole := ""
		detailsId := int32(0)
		detailsName := ""
		if smartCard.SantriID.Valid {
			ownerRole = "santri"
			detailsId = smartCard.SantriID.Int32
			detailsName = smartCard.SantriName.String
		} else if smartCard.EmployeeID.Valid {
			ownerRole = "employee"
			detailsId = smartCard.EmployeeID.Int32
			detailsName = smartCard.EmployeeName.String
		}

		result = append(result, model.SmartCardComplete{
			SmartCard: model.SmartCard{
				ID:        smartCard.ID,
				Uid:       smartCard.Uid,
				CreatedAt: smartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
				IsActive:  smartCard.IsActive,
			},
			Owner: model.OwenerDetails{
				ID:   detailsId,
				Role: repo.RoleType(ownerRole),
				Name: detailsName,
			},
		})
	}

	return &result, nil
}

func (c *service) CountSmartCards(ctx context.Context, request *model.ListSmartCardRequest) (int64, error) {
	count, err := c.store.CountSmartCards(ctx, repo.CountSmartCardsParams{
		Q:         pgtype.Text{String: request.Q, Valid: request.Q != ""},
		IsActive:  pgtype.Bool{Bool: true, Valid: true},
		CardOwner: repo.NullCardOwner{CardOwner: request.CardOwner, Valid: request.CardOwner != ""},
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *service) GetSmartCard(ctx context.Context, request *model.SmartCardRequest) (*model.SmartCardComplete, error) {
	smartCard, err := c.store.GetSmartCard(ctx, request.Uid)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("Smart Card not found")
		}
		return nil, err
	}

	var ownerRole string
	var ownerId int32
	var ownerName string

	if smartCard.SantriID.Valid {
		ownerRole = "santri"
		ownerId = smartCard.SantriID.Int32
		ownerName = smartCard.SantriName.String
	} else if smartCard.EmployeeID.Valid {
		ownerId = smartCard.EmployeeID.Int32
		ownerName = smartCard.EmployeeName.String

	} else {
		ownerRole = ""
		ownerId = 0
		ownerName = ""
	}

	return &model.SmartCardComplete{
		SmartCard: model.SmartCard{
			ID:        smartCard.ID,
			Uid:       smartCard.Uid,
			CreatedAt: smartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			IsActive:  smartCard.IsActive,
		},
		Owner: model.OwenerDetails{
			ID:   ownerId,
			Role: repo.RoleType(ownerRole),
			Name: ownerName,
		},
	}, nil
}

func (c *service) UpdateSmartCard(ctx context.Context, request *model.UpdateSmartCardRequest, id int32) (*model.SmartCardComplete, error) {
	var ownerRole string
	var detailsId int32
	var detailsName string

	if request.OwnerRole == repo.RoleTypeSantri {
		santri, err := c.store.GetSantri(ctx, request.OwnerID)
		if err != nil {
			return nil, err
		}
		detailsId = santri.ID
		detailsName = santri.Name
		ownerRole = "santri"
	} else {
		employee, err := c.store.GetEmployee(ctx, request.OwnerID)
		if err != nil {
			return nil, err
		}
		detailsId = employee.ID
		detailsName = employee.Name
		ownerRole = "employee"
	}

	updatedSmartCard, err := c.store.UpdateSmartCard(ctx, repo.UpdateSmartCardParams{
		ID:         id,
		IsActive:   pgtype.Bool{Bool: request.IsActive, Valid: true},
		SantriID:   pgtype.Int4{Int32: request.OwnerID, Valid: ownerRole == "santri"},
		EmployeeID: pgtype.Int4{Int32: request.OwnerID, Valid: ownerRole == "employee"},
	})

	if err != nil {
		return nil, err
	}

	return &model.SmartCardComplete{
		SmartCard: model.SmartCard{
			ID:        updatedSmartCard.ID,
			Uid:       updatedSmartCard.Uid,
			CreatedAt: updatedSmartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			IsActive:  updatedSmartCard.IsActive,
		},
		Owner: model.OwenerDetails{
			ID:   detailsId,
			Role: repo.RoleType(ownerRole),
			Name: detailsName,
		},
	}, nil
}

func (c *service) DeleteSmartCard(ctx context.Context, id int32) (*model.SmartCard, error) {
	deletedSmartCard, err := c.store.DeleteSmartCard(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.SmartCard{
		ID:        deletedSmartCard.ID,
		Uid:       deletedSmartCard.Uid,
		CreatedAt: deletedSmartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		IsActive:  deletedSmartCard.IsActive,
	}, nil
}
